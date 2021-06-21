package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// ErrInvalidProductName return when category name is invalid
	ErrInvalidProductName = fmt.Errorf("invalid product name")
)

// Product - requet model for variant
type Product struct {
	Name        string `json:"name"`
	CategoryID  int    `json:"category_id"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// NewAddProduct - parses request for add variants
func NewAddProduct(r *http.Request) (*Product, error) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return nil, ErrInvalidJSON
	}
	err = product.validateProduct()
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (product *Product) validateProduct() error {
	if len(product.Name) == 0 {
		return ErrInvalidProductName
	}
	if product.CategoryID <= 0 {
		return ErrInvalidCategoryID
	}
	return nil
}
