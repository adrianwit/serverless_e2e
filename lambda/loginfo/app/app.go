package main

import (
	"context"
	"encoding/json"
	"loginfo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"net/url"
	"os"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func handleRequest(ctx context.Context, apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request, err := newRequest(apiRequest)
	if err != nil {
		return handleError(err)
	}
	service := loginfo.New()
	response := service.CountLogs(request)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return handleError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

func newRequest(apiRequest events.APIGatewayProxyRequest) (request *loginfo.Request, err error) {
	request = &loginfo.Request{}
	if apiRequest.HTTPMethod == "POST" {
		return request, json.Unmarshal([]byte(apiRequest.Body), request)
	} else {
		request.Region = apiRequest.QueryStringParameters["region"]
		request.URL, err = url.QueryUnescape(apiRequest.QueryStringParameters["url"])
	}
	return request, err
}

func handleError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Printf("unable to process request %v", err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}