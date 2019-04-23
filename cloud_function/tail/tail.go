package tail

import (
	"context"
	_ "github.com/viant/bgc"
	_ "github.com/viant/toolbox/storage/gs"
)

const configKey = "config"

// TailFn loads data to matched destination
func TailFn(ctx context.Context, event GCSWorkflowEvent) (err error) {
	config, err := NewConfigFromEnv(configKey)
	if err != nil {
		return err
	}
	service, err := New(config)
	if err != nil {
		return err
	}
	return service.Transfer(ctx, event.URL())
}
