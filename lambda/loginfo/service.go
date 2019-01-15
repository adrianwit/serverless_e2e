package loginfo

import (
	"github.com/viant/toolbox/cred"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/storage/s3"
)

//Service represents a counting service
type Service interface {
	CountLogs(request *Request) *Response
}

type service struct{}

func (s *service) CountLogs(request *Request) *Response {
	response := NewResponse()
	s3.SetProvider(&cred.Config{Region: request.Region})
	service, err := storage.NewServiceForURL(request.URL, "")
	if response.SetError(err) {
		return response
	}
	countFileAndLines(service, request.URL, response)
	return response
}

func New() Service {
	return &service{}
}
