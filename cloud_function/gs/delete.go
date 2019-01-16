package gs

import (
	"cloud.google.com/go/storage"
	"context"
	"net/url"
	"os"
)

//Delete deletes gs object
func Delete(ctx context.Context, URL string) error {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "file" {
		err = os.Remove(parsedURL.Path)
		return err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bucket := client.Bucket(parsedURL.Host)
	objectPath := string(parsedURL.Path[1:])
	return bucket.Object(objectPath).Delete(ctx)
}
