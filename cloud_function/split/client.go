package split

import (
	"cloud.google.com/go/storage"
	"context"
)

var client *storage.Client

func getClient() (*storage.Client, error) {
	if client != nil {
		return client, nil
	}
	var err error
	client, err = storage.NewClient(context.Background())
	return client, err
}
