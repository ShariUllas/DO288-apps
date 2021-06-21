package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var (
	// ErrInvalidCategory return when category name is invalid
	ErrInvalidCategory = fmt.Errorf("invalid category name")
	// ErrInvalidCategoryID returns when category id is invalid
	ErrInvalidCategoryID = fmt.Errorf("invalid category id")
	// ErrParentCategoryID returns when parent category id is invalid
	ErrParentCategoryID = fmt.Errorf("invalid parent category id")
)

// Category - requet model for variant
type Category struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}

// NewAddCategory - parses request for add variants
func NewAddCategory(r *http.Request) (*Category, error) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		return nil, ErrInvalidJSON
	}
	err = category.validateCategory()
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (category *Category) validateCategory() error {
	if len(category.Name) == 0 {
		return ErrInvalidCategory
	}
	return nil
}

// NewEditCategory transforms request into edit category request
func NewEditCategory(r *http.Request) (*Category, error) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "category_id"))
	if err != nil {
		return nil, ErrInvalidCategoryID
	}

	var category Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}
	category.ID = categoryID

	return &category, nil
}

// NewDeleteCategory transforms request into category Delete
func NewDeleteCategory(r *http.Request) (*int, error) {

	categoryID, err := strconv.Atoi(chi.URLParam(r, "category_id"))
	if err != nil {
		return nil, ErrInvalidCategoryID
	}

	return &categoryID, nil
}

func NewListCategory(r *http.Request) (*int, error) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "category_id"))
	if err != nil {
		return nil, ErrInvalidCategoryID
	}

	return &categoryID, nil
}
