package model

import "time"

type DocumentType string

type LineItem struct {
	ProductID int
	Price float64
	Quantity float64
}


type Document struct {
	ID int
	Items []*LineItem
	Type DocumentType
	Created *time.Time
	PostDate *time.Time
}
