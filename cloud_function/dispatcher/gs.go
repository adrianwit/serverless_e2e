package dispatcher

import (
	"context"
	"dispatcher/gs"
	"fmt"
	"os"
)

var storageDispatcher Service

//BQEventDispatcher dispatch Storage events to matched targets
func GSEventDispatcherFn(ctx context.Context, event *gs.Event) error {
	service, err := getGSService(ctx)
	if err != nil {
		return err
	}
	return service.Handle(ctx, event)
}

func getGSService(ctx context.Context) (Service, error) {
	if storageDispatcher != nil {
		return storageDispatcher, nil
	}
	configURL := os.Getenv("configURL")
	if configURL == "" {
		return nil, fmt.Errorf("configURL was empty")
	}
	config, err := NewConfigFromURL(ctx, configURL)
	if err != nil {
		return nil, err
	}
	storageDispatcher, err = New(config, gs.NewPredicate)
	return storageDispatcher, err
}
