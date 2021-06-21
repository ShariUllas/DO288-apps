package response

import "catalogue/api/internal/core/data"

// AddProduct is
type AddProduct struct {
	ID int `json:"id"`
}

// Product - model for list product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
}

// NewAddProduct -
func NewAddProduct(id *int) AddProduct {
	return AddProduct{ID: *id}
}

// NewListProductResponse - returns response for list product
func NewListProductResponse(data []*data.Product) []*Product {
	var products []*Product
	for _, v := range data {
		var product Product
		product.ID = v.ID
		product.Name = v.Name
		product.Category = v.Category
		product.Description = v.Description
		product.ImageURL = v.ImageURL
		products = append(products, &product)
	}
	return products
}
