package dispatcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/api/pubsub/v1"
)

func Publish(ctx context.Context, target *Target, event interface{}) error {
	httpClient, err := getDefaultClient(ctx, pubsub.CloudPlatformScope, pubsub.PubsubScope)
	if err != nil {
		return err
	}
	service, err := pubsub.New(httpClient)
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
