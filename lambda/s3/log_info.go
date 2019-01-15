package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stderr, "INFO ", log.Llongfile)

type book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	JSON, err := json.Marshal(getBook())
	if err != nil {
		return serverError(err)
	}
	body := string(JSON)

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func getBook() *book {
	bk := &book{
		ISBN:   "978-1420931693",
		Title:  "The Republic",
		Author: "Plato",
	}

	return bk
}

func main() {
	lambda.Start(handle)
}
