## Lambda with e2e testing

**Prerequisites:**

 - go1.11
 - dedicated aws account for e2e testing 
    * aws credentials file-> ~/.secret/e2e.json 
    * [Setup Endly AWS Credentials](https://github.com/viant/endly/tree/master/doc/secrets#aws)
    
#### Introduction

AWS Lambda is a serverless compute service that runs a code in response to events and automatically manages the underlying compute resources.
The event can be fired by a specific trigger, which determines how and when your function executes. 


Lambda entry point:
```go
func main() {
	lambda.Start(handler)
}
``` 

where handler supports the following function signatures: 

- func()
- func() error
- func(Event) error
- func() (Output, error)
- func(Event) (Output, error)
- func(context.Context) error
- func(context.Context, Event) error
- func(context.Context) (Output, error)
- func(context.Context, Event) (Output, error)


This project provides example for the following native mechanisms:

#### Direct lambda trigger

- [HelloWorld](hello/hello.go)
- [E2E Use Case](e2e/regression/cases/001_hello_world)
- Signature 
```go
    func() (Output, error)
```

#### HTTP APIGateway trigger


**APIGateway Proxy Event**
```go
// APIGatewayProxyRequest contains data coming from the API Gateway proxy
type APIGatewayProxyRequest struct {
	Resource                        string                        `json:"resource"` // The resource path defined in API Gateway
	Path                            string                        `json:"path"`     // The url path for the caller
	HTTPMethod                      string                        `json:"httpMethod"`
	Headers                         map[string]string             `json:"headers"`
	MultiValueHeaders               map[string][]string           `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string             `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string           `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string             `json:"pathParameters"`
	StageVariables                  map[string]string             `json:"stageVariables"`
	RequestContext                  APIGatewayProxyRequestContext `json:"requestContext"`
	Body                            string                        `json:"body"`
	IsBase64Encoded                 bool                          `json:"isBase64Encoded,omitempty"`
}
```

- [LogInfo](loginfo/app/loginfo.go)
- [E2E Use Case](e2e/regression/cases/002_logs_count)
- Signature 
```go
    func(ctx context.Context, apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
```
- Resources
    * [Amazon API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html)
    * [Event Contract](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go)
- Usage    
    
```go

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request, err := newRequest(apiRequest)
	if err != nil {
		return handleError(err)
	}
	service := New() //you service
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


func handleError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Printf("unable to process request %v", err)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

```    
