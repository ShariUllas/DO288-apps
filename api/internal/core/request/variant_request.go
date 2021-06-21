package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var (
	// ErrInvalidJSON retuns when json is invalid
	ErrInvalidJSON = fmt.Errorf("invalid JSON received")
	// ErrInvalidMRP retuns when MRP is invalid
	ErrInvalidMRP = fmt.Errorf("invalid MRP value received")
	// ErrInvalidVariantID returns when variant id is invalid
	ErrInvalidVariantID = fmt.Errorf("invalid variant id received")
	// ErrInvalidProductID returns when product id is invalid
	ErrInvalidProductID = fmt.Errorf("invalid product id received")
)

// Variant - requet model for variant
type Variant struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	MRP           int    `json:"mrp"`
	DiscountPrice int    `json:"discount"`
	Size          int    `json:"size"`
	Color         string `json:"color"`
	ProductID     int    `json:"product_id,required"`
}

// NewAddVariant - parses request for add variants
func NewAddVariant(r *http.Request) (*Variant, error) {
	var variant Variant
	err := json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		return nil, ErrInvalidJSON
	}
	validate(&variant)
	return &variant, nil
}

func validate(variant *Variant) error {
	if variant.MRP == 0 {
		return ErrInvalidMRP
	}
	return nil
}

// NewEditVariant transforms request into edit variant request
func NewEditVariant(r *http.Request) (*Variant, error) {
	variantID, err := strconv.Atoi(chi.URLParam(r, "variant_id"))
	if err != nil {
		return nil, ErrInvalidVariantID
	}

	var variant Variant
	err = json.NewDecoder(r.Body).Decode(&variant)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	variant.ID = variantID
	if variant.ProductID <= 0 {
		return nil, ErrInvalidProductID
	}
	if variant.MRP <= 0 {
		return nil, ErrInvalidMRP
	}
	return &variant, nil
}

// NewDeleteVariant transforms request into UserDelete
func NewDeleteVariant(r *http.Request) (*int, error) {

	variantID, err := strconv.Atoi(chi.URLParam(r, "variant_id"))
	if err != nil {
		return nil, ErrInvalidVariantID
	}

	return &variantID, nil
}
