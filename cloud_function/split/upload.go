package split

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"net/url"
)

//Upload uploads content to google storage
func Upload(ctx context.Context, URL string, reader io.Reader) error {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bucket := client.Bucket(parsedURL.Host)
	objectPath := string(parsedURL.Path[1:])
	writer := bucket.Object(objectPath).NewWriter(ctx)
	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}
	return writer.Close()
}
