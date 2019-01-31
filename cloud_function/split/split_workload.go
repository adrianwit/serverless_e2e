package split

import (
	"context"
	"fmt"
	"github.com/adrianwit/serverless_e2e/cloud_function/gs"
	"github.com/viant/toolbox"
	"path"
	"strings"
)

const (
	fragmentCount = 3
	workerPath    = "data/workers/"
	masterPath    = "data/master/"
)

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

// SplitWorkloadFn split incoming master files into a smaller workers files
func SplitWorkloadFn(ctx context.Context, event GCSWorkflowEvent) error {
	URL := fmt.Sprintf("gs://%v/%v", event.Bucket, event.Name)
	payload, err := gs.Download(ctx, URL)
	if err != nil {
		return err
	}
	if !strings.Contains(event.Name, masterPath) {
		return nil
	}
	fragments := toolbox.TerminatedSplitN(string(payload), fragmentCount, "\n")
	_, ownerName := path.Split(event.Name)
	for i, fragment := range fragments {
		fragmentURL := fmt.Sprintf("gs://%v/%v", event.Bucket, path.Join(workerPath, fmt.Sprintf("%02d_%s", i, ownerName)))
		if err = gs.Upload(ctx, fragmentURL, strings.NewReader(fragment)); err != nil {
			return fmt.Errorf("failed to upload %v %v", fragmentURL, err)
		}
	}
	return nil
}
