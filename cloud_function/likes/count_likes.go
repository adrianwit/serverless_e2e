package cloud_function

import (
	"cloud.google.com/go/functions/metadata"
	"context"
)

//FirebaseEvent represents a firebase event
type FirebaseEvent struct {
	Delta interface{} `json:"delta"`
	Auth  interface{} `json:"auth"`
}

var srv Service

func init() {
	srv = New()
}

//CountLikesFn counts post likes to update likes_count
func CountLikesFn(ctx context.Context, event FirebaseEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return err
	}
	resource := Resource(meta.Resource.Name)
	return srv.UpdateLikes(ctx, resource)
}
