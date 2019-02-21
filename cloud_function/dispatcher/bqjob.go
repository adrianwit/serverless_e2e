package dispatcher

import (
	"golang.org/x/net/context"
	"google.golang.org/api/bigquery/v2"
)

var bqService *bigquery.Service

func getBQService() (*bigquery.Service, error) {
	if bqService != nil {
		return bqService, nil
	}
	httpClient, err := getDefaultClient(context.Background(), bigquery.CloudPlatformScope, bigquery.BigqueryScope, bigquery.BigqueryInsertdataScope)
	if err != nil {
		return nil, err
	}
	bqService, err = bigquery.New(httpClient)
	return bqService, err
}

func GetBQJob(ctx context.Context, projectID, jobID string) (*bigquery.Job, error) {
	service, err := getBQService()
	if err != nil {
		return nil, err
	}
	jobService := bigquery.NewJobsService(service)
	call := jobService.Get(projectID, jobID)
	call.Context(ctx)
	return call.Do()
}
