package agg

import "time"

type EventKey struct {
	Date string `json:"EventDate"`
	EventType int	`json:"EventType"`
}

type Message struct {
	*EventKey
	Timestamp *time.Time `json:"Timestamp"`
}

