package response

import "catalogue/api/internal/core/data"

// Variant is response model for variant
type Variant struct {
	ID            int     `json:"id"`
	Name          *string `json:"name"`
	MRP           int     `json:"mrp"`
	DiscountPrice *int    `json:"discount"`
	Size          *int    `json:"size"`
	Color         *string `json:"color"`
	Product       string  `json:"product"`
}

// NewVariantResponse returns response model for variant
func NewVariantResponse(data []*data.Variant) []*Variant {
	var variants []*Variant
	for _, v := range data {
		var variant Variant
		variant.ID = v.ID
		variant.Color = v.Color
		variant.DiscountPrice = v.DiscountPrice
		variant.MRP = v.MRP
		variant.Name = v.Name
		variant.Size = v.Size
		variant.Product = v.Product
		variants = append(variants, &variant)
	}
	return variants
}
