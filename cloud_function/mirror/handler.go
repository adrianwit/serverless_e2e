package split

import (
	"context"
	"fmt"
	"runtime/debug"
)



const configKey = "config"

// GCSWorkflowEvent is the payload of a GCS event.
type GCSWorkflowEvent struct {
	Bucket      string `json:"bucket"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	CRC32C      string `json:"crc32c"`
	Kind        string `json:"kind"`
	Size        string `json:"size"`
	SelfLink    string `json:"selfLink"`
	MediaLink   string `json:"mediaLink"`
}



// MirrorFn copies source to destination
func MirrorFn(ctx context.Context, event GCSWorkflowEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = fmt.Errorf("%v", r)
		}
	}()
	URL := fmt.Sprintf("gs://%v/%v", event.Bucket, event.Name)
	config, err := NewConfigFromEnv(configKey)
	if err != nil {
		return err
	}
	service, err := GetService(ctx, config)
	if err != nil {
		return err
	}
	return service.Mirror(URL)
}
