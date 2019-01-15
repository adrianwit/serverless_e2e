package loginfo

import "fmt"

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type Request struct {
	Region string
	URL    string
}

func (r *Request) Validate() error {
	if r.Region == "" {
		return fmt.Errorf("region was empty")
	}
	if r.URL == "" {
		return fmt.Errorf("url was empty")
	}
	return nil
}


type Response struct {
	Status     string
	Error      string
	FileCount  uint32
	LinesCount uint32
}


func (r *Response) SetError(err error) bool {
	if err == nil {
		return false
	}
	r.Status = StatusError
	r.Error = err.Error()
	return true
}

func NewResponse() *Response {
	return &Response{
		Status:StatusOK,
	}
}

