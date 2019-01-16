package fs

import (
	"cloud.google.com/go/firestore"
	"context"
)

//RunTransaction runs transactions
func RunTransaction(ctx context.Context, projectID, path, id string, updateFunc func(doc *firestore.DocumentRef, transaction *firestore.Transaction) error) error {
	client, err := NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	ref := client.Collection(path).Doc(id)
	return client.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {
		return updateFunc(ref, transaction)
	})
}
