package tail

import "fmt"

// GCSWorkflowEvent is the payload of a GCS event.
type GCSWorkflowEvent struct {
	Bucket      string `json:"bucket"`
	Name        string `json:"name"`
}

func (e *GCSWorkflowEvent) URL() string {
	return fmt.Sprintf("gs://%v/%v", e.Bucket, e.Name)
}

