package msg

import (
	"cloud.google.com/go/pubsub"
	"context"
)

func Publish(ctx context.Context, projectID string, topicID string, messages ...*pubsub.Message) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	topic := client.Topic(topicID)
	for _, message := range messages {
		response := topic.Publish(ctx, message)
		if _, err = response.Get(ctx); err != nil {
			return err
		}
	}
	return nil
}
