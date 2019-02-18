package computeinfo

import (
	"cloud.google.com/go/functions/metadata"
	"context"
	"encoding/json"
	"fmt"
)

type Event struct {
	Name string `json:"name"`
}

func OnChangeFn(ctx context.Context, event *Event) error {
	fmt.Printf("instance %v is down\n", event.Name)
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %v", err)
	}

	JSON, _ := json.Marshal(meta)
	fmt.Printf("meta: %v\n", string(JSON))
	return nil
}
