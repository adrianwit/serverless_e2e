## Lambda with e2e testing

**Prerequisites:**

 - go1.11
 - dedicated aws account for e2e testing 
    * aws credentials file-> ~/.secret/e2e.json 
    * [Setup Endly AWS Credentials](https://github.com/viant/endly/tree/master/doc/secrets#aws)
    
#### Introduction

AWS Lambda is a serverless compute service that runs a provided code in response to events and automatically manages the underlying compute resources.
The event can be fired by a specific trigger, which determines how and when your function executes. 

References:
* [Programming-model](https://docs.aws.amazon.com/lambda/latest/dg/programming-model-v2.html)

    
#### Event/Trigger

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

References:
* [Go Programming Model Handler Types](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html)


### Examples

#### Direct lambda trigger

- [HelloWorld](hello/hello.go)
- [E2E Use Case](e2e/regression/cases/001_hello_world)
- _Handler Signature_ 
```go
    func() (Output, error)
    func(Event) (Output, error)
```


#### HTTP APIGateway trigger

- [LogInfo](loginfo/app/loginfo.go)
- [E2E Use Case](e2e/regression/cases/002_logs_count)
- _Handler Signature_ 
```go
    func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
```
- APIGatewayProxyRequest 
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
- References
    * [Amazon API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html)
    * [Event Contract](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go)


#### S3 triggers

- [FileMeta](filemeta/filemeta.go)
- [E2E Use Case](e2e/regression/cases/003_filemeta)
- _Handler Signature_ 
```go
    func(context.Context, events.S3Event)
```

- Event types:
 * s3:ObjectCreated:*
 * s3:ObjectCreated:Put
 * s3:ObjectCreated:Post
 * s3:ObjectCreated:Copy
 * s3:ObjectCreated:CompleteMultipartUpload
 * s3:ObjectRemoved:*
 * s3:ObjectRemoved:Delete
 * s3:ObjectRemoved:DeleteMarkerCreated
 * s3:ObjectRestore:Post
 * s3:ObjectRestore:Completed
 * s3:ReducedRedundancyLostObject 

- S3Event
```go
type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

type S3EventRecord struct {
	EventVersion      string              `json:"eventVersion"`
	EventSource       string              `json:"eventSource"`
	AWSRegion         string              `json:"awsRegion"`
	EventTime         time.Time           `json:"eventTime"`
	EventName         string              `json:"eventName"`
	PrincipalID       S3UserIdentity      `json:"userIdentity"`
	RequestParameters S3RequestParameters `json:"requestParameters"`
	ResponseElements  map[string]string   `json:"responseElements"`
	S3                S3Entity            `json:"s3"`
}

type S3UserIdentity struct {
	PrincipalID string `json:"principalId"`
}
```
- [Event Contract](https://github.com/aws/aws-lambda-go/blob/master/events/s3.go)

- References:
    * [Using AWS Lambda with Amazon S3](https://docs.aws.amazon.com/lambda/latest/dg/with-s3.html)
    * [Using AWS Lambda with Amazon S3 Example](https://docs.aws.amazon.com/lambda/latest/dg/with-s3-example.html)
    * [NotificationHowTo](https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html)
    * [Setting Bucket Notification](https://docs.aws.amazon.com/cli/latest/reference/s3api/put-bucket-notification-configuration.html)



### Error Handling

References:
 * [Errors](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-errors.html)
 