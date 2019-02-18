package dispatcher

import (
	"cloud.google.com/go/functions/metadata"
	"context"
	"dispatcher/bq"
	"fmt"
	"os"
	"strings"
)

var bqSrv Service

//BQEventDispatcher dispatches BigQuery events to matched targets
func BQEventDispatcherFn(ctx context.Context, eventData struct{}) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %v", err)
	}
	resourceParts := strings.Split(meta.Resource.Name, "/")
	projectID := resourceParts[1]
	jobID := resourceParts[len(resourceParts)-1]
	job, err := GetBQJob(ctx, projectID, jobID)
	if err != nil {
		return err
	}
	event := bq.Event(*job)
	service, err := getService(ctx)
	if err != nil {
		return err
	}
	err = service.Handle(ctx, &event)
	return err
}

func getService(ctx context.Context) (Service, error) {
	if bqSrv != nil {
		return bqSrv, nil
	}
	configURL := os.Getenv("configURL")
	if configURL == "" {
		return nil, fmt.Errorf("configURL was empty")
	}
	config, err := NewConfigFromURL(ctx, configURL)
	if err != nil {
		return nil, err
	}
	bqSrv, err = New(config, bq.NewPredicate)
	return bqSrv, err
}
