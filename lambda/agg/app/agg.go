package main

import (
	"agg"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func handleEvent(ctx context.Context, snsEvent events.SNSEvent) {
	service := agg.New()
	for _, record := range snsEvent.Records {
		message, err := agg.SNSEventRecordToMessage(&record)
		if err == nil {
			err = service.Consume(message)
		}
		if err != nil {
			errorLogger.Printf("%v", err)
			continue
		}
	}
}

func main() {
	lambda.Start(handleEvent)
}
