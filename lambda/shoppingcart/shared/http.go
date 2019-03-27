package shared

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
request := &http.Request{
		Method:apiRequest.HTTPMethod,
		URL:
	}
 */



 func NewHttpRequest(apiRequest *events.APIGatewayProxyRequest) (*http.Request, error) {
	request := &http.Request{
		Method:apiRequest.HTTPMethod,
		Header:http.Header{},
	}
	URL, err := url.Parse(fmt.Sprintf("https://localhost/%v", apiRequest.Path))
	if err != nil {
		return nil, err
	}
	request.URL = URL

	if apiRequest.Body != "" {
		if apiRequest.IsBase64Encoded {
			payload, err := base64.StdEncoding.DecodeString(apiRequest.Body)
			if err != nil {
				return nil, err
			}
			request.Body = ioutil.NopCloser(bytes.NewReader(payload))
		} else {
			request.Body = ioutil.NopCloser(strings.NewReader(apiRequest.Body))
		}
	}
	for k, v := range  apiRequest.MultiValueHeaders {
		request.Header[k]= v
	}
	return request, nil
}


func SetCORSHeaderIfNeeded(apiRequest *events.APIGatewayProxyRequest, response *events.APIGatewayProxyResponse) {
	origin, ok := apiRequest.Headers["Origin"]
	if ! ok {
		return
	}
	if len(response.Headers) == 0 {
		response.Headers = make(map[string]string)
	}
	response.Headers["Access-Control-Allow-Credentials"] = "true"
	response.Headers["Access-Control-Allow-Origin"] = origin
	response.Headers["Access-Control-Allow-Methods"] = "POST GET"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type, *"
	response.Headers["Access-Control-Max-Age"] = "120"
}
