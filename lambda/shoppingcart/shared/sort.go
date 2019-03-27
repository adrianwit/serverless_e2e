package shared

import "shoppingcart/model"

type ProductSorter struct {
	products []*model.Product
	by func(p1, p2 *model.Product) bool
}



// Len is part of sort.Interface.
func (s *ProductSorter) Len() int {
	return len(s.products)
}

// Swap is part of sort.Interface.
func (s *ProductSorter) Swap(i, j int) {
	s.products[i], s.products[j] = s.products[j], s.products[i]
}
// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ProductSorter) Less(i, j int) bool {
	return s.by(s.products[i], s.products[j])
}

//NewProductSorter returns new product sorter
func NewProductSorter(products []*model.Product, by func(p1, p2 *model.Product) bool) *ProductSorter {
	return &ProductSorter{products:products, by:by}
}

//ByProductID ID comparator
func ByProductID(p1, p2 *model.Product) bool {
	return p1.ID < p2.ID
}