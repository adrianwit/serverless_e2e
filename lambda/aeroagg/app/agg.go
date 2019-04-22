package main

import (
	"aeroagg"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const envConfigKey = "CONFIG"


func handleEvent(ctx context.Context, snsEvent events.SNSEvent) error {
	service, err := getService()
	if err != nil {
		return err
	}

	for _, record := range snsEvent.Records {
		message, err := aeroagg.SNSEventRecordToMessage(&record)
		if err == nil {
			err = service.Consume(message)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

var service aeroagg.Service

func getService() (aeroagg.Service, error) {
	if service != nil {
		return service, nil
	}
	config, err := aeroagg.NewConfigFromEnv(envConfigKey)
	if err != nil {
		return nil, err
	}
	service = aeroagg.New(config)
	return service, nil
}

func main() {
	lambda.Start(handleEvent)
}
