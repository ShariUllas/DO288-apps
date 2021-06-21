package response

import "catalogue/api/internal/core/data"

// AddCategory is
type AddCategory struct {
	ID int `json:"id"`
}

// Category -
type Category struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category *string `json:"category"`
}
type SubCategory struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	MainCategory []*Category `json:"sub_categories"`
}

// NewAddCategory -
func NewAddCategory(id *int) AddCategory {
	return AddCategory{ID: *id}
}

// NewListCategoryResponse -
func NewListCategoryResponse(data []*data.Category) []*Category {
	var categories []*Category
	for _, v := range data {
		var category Category
		category.ID = v.ID
		category.Name = v.Name
		category.Category = v.Category
		categories = append(categories, &category)
	}
	return categories
}

func NewCategories(main *data.Category, data []*data.Category) *SubCategory {
	var categories []*Category
	for _, v := range data {
		var category Category
		category.ID = v.ID
		category.Name = v.Name
		categories = append(categories, &category)
	}
	var subcategory SubCategory
	subcategory.ID = main.ID
	subcategory.Name = main.Name
	subcategory.MainCategory = categories
	return &subcategory
}
