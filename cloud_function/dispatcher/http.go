package dispatcher

import (
	"context"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
	"net/http"
)

func getDefaultClient(ctx context.Context, scopes ...string) (*http.Client, error) {
	o := []option.ClientOption{
		option.WithScopes(scopes...),
	}
	httpClient, _, err := htransport.NewClient(ctx, o...)
	return httpClient, err
}
