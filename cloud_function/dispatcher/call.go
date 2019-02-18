package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/storage/v1"
	"net/http"
	"net/url"
	"os"
)

const functionURL = "https://%s-%s.cloudfunctions.net/%s"

func Call(ctx context.Context, target *Target, event interface{}) error {
	httpClient, err := getDefaultClient(ctx, cloudfunctions.CloudPlatformScope, storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	URL := getURL(target)
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("post", URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	_, err = httpClient.Do(request)
	return err
}

func getURL(target *Target) string {
	if _, err := url.Parse(target.URL); err == nil {
		return target.URL
	}
	projectID := os.Getenv("GCLOUD_PROJECT")
	region := os.Getenv("FUNCTION_REGION")
	return fmt.Sprintf(functionURL, region, projectID, target.URL)
}
