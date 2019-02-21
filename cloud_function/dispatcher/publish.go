package dispatcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/api/pubsub/v1"
)

var pubsubService *pubsub.Service

func getPubsubService() (*pubsub.Service, error) {
	if pubsubService != nil {
		return pubsubService, nil
	}
	httpClient, err := getDefaultClient(context.Background(), pubsub.CloudPlatformScope, pubsub.PubsubScope)
	if err != nil {
		return nil, err
	}
	pubsubService, err = pubsub.New(httpClient)
	return pubsubService, err
}

func Publish(ctx context.Context, target *Target, event interface{}) error {
	service, err := getPubsubService()
	if err != nil {
		return err
	}
	topicService := pubsub.NewProjectsTopicsService(service)
	message := &pubsub.PubsubMessage{}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal even: %v", err)
	}
	message.Data = base64.StdEncoding.EncodeToString(data)
	request := &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			message,
		},
	}
	call := topicService.Publish(target.URL, request)
	call.Context(ctx)
	_, err = call.Do()
	return err
}
