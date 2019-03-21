package main

import (
	"context"
	"dstransfer"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

const configKey = "CONFIG"

var service dstransfer.Service

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func main() {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Println("Recovered in f", r)
		}
	}()
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request, err := newRequest(apiRequest)
	if err != nil {
		return handleError(err)
	}
	service, err := getService()
	if err != nil {
		return handleError(err)
	}

	response := service.Copy(request)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return handleError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}

func newRequest(apiRequest events.APIGatewayProxyRequest) (request *dstransfer.Request, err error) {
	request = &dstransfer.Request{}
	if apiRequest.HTTPMethod == "POST" {
		return request, json.Unmarshal([]byte(apiRequest.Body), request)
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

func handleEvent(request *dstransfer.Request) (*dstransfer.Response, error) {
	service, err := getService()
	if err != nil {
		return nil, err
	}
	response := service.Copy(request)
	return response, nil
}

func getService() (dstransfer.Service, error) {
	if service != nil {
		return service, nil
	}
	config := &dstransfer.Config{}
	err := json.NewDecoder(strings.NewReader(os.Getenv(configKey))).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}
	service, err = dstransfer.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to crate service: %v", err)
	}
	return service, err
}
