package cloud_function

import (
	"cloud.google.com/go/functions/metadata"
	"context"
	"firebase.google.com/go/db"
	"fmt"
	"github.com/adrianwit/serverless_e2e/cloud_function/fb"
	"strings"
)

//FirebaseEvent represents a firebase event
type FirebaseEvent struct {
	Delta interface{} `json:"delta"`
	Auth  interface{} `json:"auth"`
}

/*
gcloud alpha functions deploy CountLikesFn
--trigger-event providers/google.firebase.database/eventTypes/ref.write
--trigger-resource  'projects/_/instances/$instanceID/refs/posts/{key}/likes'
--runtime go111
*/

//CountLikesFn counts post likes to update likes_count
func CountLikesFn(ctx context.Context, event FirebaseEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return err
	}
	fragments := strings.Split(meta.Resource.Name, "/")
	instanceID := fragments[3]
	key := fragments[len(fragments)-2]
	databaseURL := fmt.Sprintf("https://%s.firebaseio.com", instanceID)
	refPath := fmt.Sprintf("posts/%v", key)
	return fb.RunTransaction(ctx, databaseURL, refPath, func(nodeSnapshot db.TransactionNode) (interface{}, error) {
		var record = make(map[string]interface{})
		if err := nodeSnapshot.Unmarshal(&record); err != nil {
			return nil, err
		}
		if likes, ok := record["likes"]; ok {
			if likesMap, ok := likes.(map[string]interface{}); ok {
				record["likes_count"] = len(likesMap)
			}
		}
		return record, nil
	})
}
