package split

import (
	"cloud.google.com/go/storage"
	"context"
	"io/ioutil"
	"net/url"
)

//Download downloads google storage content
func Download(ctx context.Context, URL string) ([]byte, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(parsedURL.Host)
	objectPath := string(parsedURL.Path[1:])
	rc, err := bucket.Object(objectPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}
