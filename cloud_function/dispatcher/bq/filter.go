package bq

//Filter represents big query job filter
type Filter struct {
	JobID       string
	Source      string
	Destination string
	Type        string
	Status      string
}
