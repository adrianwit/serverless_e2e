package main

import (
	"context"
	"filemeta"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"

)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
const metaURL = "meta/filemeta.json"

func handleEvent(ctx context.Context, s3Event events.S3Event) {
	if len(s3Event.Records) ==  0 {
		return
	}
	service, err := filemeta.New(s3Event.Records[0].AWSRegion)
	if err != nil {
		errorLogger.Printf("unable to create storage service %v\n", err)
		return
	}
	request := &filemeta.Request{
		MetaURL: fmt.Sprintf("s3://%s/%s",  s3Event.Records[0].S3.Bucket.Name, metaURL),
		ObjectURLs: make([]string, 0),
	}

	for _, record := range s3Event.Records {
		URL := fmt.Sprintf("s3://%s/%s", record.S3.Bucket.Name, record.S3.Object.Key)
		request.ObjectURLs = append(request.ObjectURLs, URL)
	}
	err = service.UpdateMeta(request)
	if err != nil {
		errorLogger.Printf("unable to update meta %v\n", err)
		return
	}
}



func main() {
	lambda.Start(handleEvent)
}
