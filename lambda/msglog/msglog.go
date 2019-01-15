package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin/json"
)

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("Body: %v", message.Body)
		JSONMesage, _ := json.Marshal(message)
		fmt.Printf("%s", JSONMesage)
	}
	return nil
}



func main() {
	lambda.Start(handleMessages)
}