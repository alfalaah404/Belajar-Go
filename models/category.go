package models

// Category represents a product category
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Elektronik"`
	Description string `json:"description" example:"Produk elektronik dan gadget"`
}
