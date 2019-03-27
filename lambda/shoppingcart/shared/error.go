package shared

import (
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"os"
)


var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)


func HandleError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Printf("unable to process request %v", err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}
