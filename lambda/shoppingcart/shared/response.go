package shared


const (
	ResponseStatusOK = "ok"
	ResponseStatusError = "error"
)

//Response represents generic response
type Response struct {
	Status string
	Error string
}


func (r *Response) SetError(err error) {
	if err == nil {
		return 
	}
	r.Status = ResponseStatusError
	r.Error = err.Error()
}


//NewResponse creates a new response
func NewResponse() Response {
	return Response{
		Status:ResponseStatusOK,
	}
}

