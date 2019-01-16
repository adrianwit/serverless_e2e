package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Printf("recived: %v\n", len(sqsEvent.Records))
	for _, message := range sqsEvent.Records {
		fmt.Printf("%v\n", message.Body)
		JSONMesage, _ := json.Marshal(message)
		fmt.Printf("%s\n", JSONMesage)
	}
	return nil
}

func main() {
	lambda.Start(handleMessages)
}
