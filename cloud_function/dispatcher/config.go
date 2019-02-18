package dispatcher

import (
	"encoding/json"
	"golang.org/x/net/context"
)

//Defines service routes
type Config struct {
	Routes []*Route
}

//NewConfigFromURL creates a new config from URL
func NewConfigFromURL(ctx context.Context, URL string) (*Config, error) {
	payload, err := Download(ctx, URL)
	if err != nil {
		return nil, err
	}
	result := &Config{}
	return result, json.Unmarshal(payload, result)
}
