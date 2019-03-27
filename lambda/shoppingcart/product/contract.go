package product

import (
	"shoppingcart/model"
	"shoppingcart/shared"
)



//ListRequest represents list product request
type ListRequest shared.ListRequest

//ListRequest represents list product response
type ListResponse struct {
	shared.Response
	Data []*model.Product
}
