package agg

type Event struct {
	EventType int `json:"EventType"`
	Date string `json:"EventDate"`
	Quantity int `json:"Qty"`
}
