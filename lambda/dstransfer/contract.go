package dstransfer

//Request represents a copy request
type Request struct {
	SQL string
	DestURL string
}

//Request represents a copy response
type Response struct {
	Status string
	Message string
	Count int
}