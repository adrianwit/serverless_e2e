package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/viant/toolbox"
	"net/url"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
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

var fileServer = http.FileServer(http.Dir("build"))

func handleRequest(ctx context.Context, proxyRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request, err := NewHttpRequest(&proxyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError) + err.Error(),
		}, nil
	}


	toolbox.DumpIndent(proxyRequest, true)

	proxyResponse := events.APIGatewayProxyResponse{
		MultiValueHeaders: make(map[string][]string),
	}
	writer := NewHttpWriter(&proxyResponse)
	fileServer.ServeHTTP(writer, request)

	toolbox.DumpIndent(proxyResponse, true)
	return proxyResponse, nil
}




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


type responseWriter struct {
	proxy *events.APIGatewayProxyResponse
}

func (r *responseWriter) Header() http.Header {
	return r.proxy.MultiValueHeaders
}


func (r *responseWriter) Write(bs []byte) (int, error){
	r.proxy.Body += string(bs)
	return len(bs), nil
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.proxy.StatusCode = statusCode
}

func NewHttpWriter(proxy *events.APIGatewayProxyResponse) http.ResponseWriter {
	return  &responseWriter{
		proxy:proxy,
	}
}

