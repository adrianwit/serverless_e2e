package cloud_function

import (
	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adrianwit/serverless_e2e/cloud_function/ps"
)

//ProxyMessage represent a proxy message
type ProxyMessage struct {
	Source    string `json:"source"`
	Dest      string `json:"dest"`
	ProjectID string `json:"projectID"`
	Message   string `json:"message"`
}

//Validate checks if message is valid
func (m *ProxyMessage) Validate() error {
	if m == nil {
		return fmt.Errorf("proxy message was nil")
	}
	if m.Dest == "" {
		return fmt.Errorf("dest was empty")
	}
	if m.ProjectID == "" {
		return fmt.Errorf("ProjectID was empty")
	}
	return nil
}

// PubsubProxyEvent
type PubsubProxyEvent struct {
	Data []byte `json:"data"`
}

// PubSubProxyFn proxies received message to specified destination  - entry point
func PubSubProxyFn(ctx context.Context, event PubsubProxyEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	message := &ProxyMessage{}
	if err = json.Unmarshal(event.Data, &message); err != nil {
		return fmt.Errorf("failed to unmarshal %s,  %v", event.Data, err)
	}
	err = message.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate message %v", err)
	}
	message.Source = meta.Resource.Name

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message %v", err)
	}
	msg := &pubsub.Message{Data: data}
	err = ps.Publish(ctx, message.ProjectID, message.Dest, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}
