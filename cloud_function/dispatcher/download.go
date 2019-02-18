package dispatcher

import (
	"context"
	"fmt"
	"google.golang.org/api/storage/v1"
	"io/ioutil"
	"net/url"
)

//Download downloads google storage content
func Download(ctx context.Context, URL string) ([]byte, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Scheme != "gs" {
		return nil, fmt.Errorf("unsupported proto scheme: %v", parsedURL.Scheme)
	}

	httpClient, err := getDefaultClient(ctx, storage.CloudPlatformScope, storage.DevstorageFullControlScope)
	service, err := storage.New(httpClient)
	if err != nil {
		return nil, err
	}
	objectService := storage.NewObjectsService(service)
	object := objectService.Get(parsedURL.Host, string(parsedURL.Path[1:]))
	object.Context(ctx)
	response, err := object.Download()
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
