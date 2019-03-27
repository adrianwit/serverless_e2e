package main

import (
	"context"
	"encoding/json"
	"fmt"
	"shoppingcart/shared"
	"shoppingcart/product"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"net/http"
	"runtime/debug"
)


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
		return shared.HandleError(err)
	}
	service := product.New()

	response := service.List(request)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return shared.HandleError(err)
	}
	apiResponse :=  events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}
	shared.SetCORSHeaderIfNeeded(&apiRequest, &apiResponse)
	return apiResponse, nil
}

func newRequest(apiRequest events.APIGatewayProxyRequest) (request *product.ListRequest, err error) {
	request = &product.ListRequest{}
	if apiRequest.HTTPMethod == "POST" {
		return request, json.Unmarshal([]byte(apiRequest.Body), request)
	}
	return request, err
}
