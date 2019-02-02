package split

import (
	"cloud.google.com/go/storage"
	"context"
	"net/url"
	"os"
)

//Exists returns true if  gs object exists
func Exists(ctx context.Context, URL string) bool {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return false
	}
	if parsedURL.Scheme == "file" {
		_, err := os.Stat(parsedURL.Path)
		return err == nil
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false
	}
	bucket := client.Bucket(parsedURL.Host)
	objectPath := string(parsedURL.Path[1:])
	_, err = bucket.Object(objectPath).Attrs(ctx)
	return err == nil
}
