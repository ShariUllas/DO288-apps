package service

import (
	"catalogue/api/internal/core/data"
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/core/response"
	"context"
	"fmt"
)

// Category is service interface
type Category interface {
	GetCategory(ctx context.Context) ([]*response.Category, error)
	AddCategory(ctx context.Context, category *request.Category) (response.AddCategory, error)
	EditCategory(ctx context.Context, category *request.Category) error
	DeleteCategory(ctx context.Context, id *int) error
	ListCategory(ctx context.Context, id *int) (*response.SubCategory, error)
}

type categoryImpl struct {
	data.CategoryRepo
}

var (
	// ErrCategoryNotFound is returned when category is not found
	ErrCategoryNotFound = fmt.Errorf("category not found")
)

// NewCategoryService is
func NewCategoryService(categoryRepo data.CategoryRepo) Category {
	return &categoryImpl{CategoryRepo: categoryRepo}
}
func (i *categoryImpl) ListCategory(ctx context.Context, id *int) (*response.SubCategory, error) {
	var sub *response.SubCategory
	maincategory, err := i.CategoryRepo.GetCategoryID(ctx, *id)
	if err != nil {
		return sub, err
	}
	res, err := i.CategoryRepo.ListCategory(ctx, id)
	if err != nil {
		return sub, err
	}
	sub = response.NewCategories(maincategory, res)
	return sub, nil
}
func (i *categoryImpl) AddCategory(ctx context.Context, category *request.Category) (response.AddCategory, error) {
	categoryID, err := i.CategoryRepo.AddCategory(ctx, category)
	if err != nil {
		return response.AddCategory{}, err
	}
	return response.NewAddCategory(categoryID), nil
}

func (i *categoryImpl) GetCategory(ctx context.Context) ([]*response.Category, error) {
	categories, err := i.CategoryRepo.GetCategory(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewListCategoryResponse(categories), nil
}

func (i *categoryImpl) EditCategory(ctx context.Context, req *request.Category) error {

	// request may contain partial data so we overwrite the existing with it
	categoryToUpdate, err := i.CategoryRepo.GetCategoryID(ctx, req.ID)
	if err != nil {
		if err == data.ErrCategoryNotFound {
			return ErrCategoryNotFound
		}
		return err
	}
	if len(req.Name) == 0 {
		categoryToUpdate.Name = req.Name
	}
	if req.ParentID <= 0 {
		categoryToUpdate.ParentID = req.ParentID
	}

	err = i.CategoryRepo.EditCategory(ctx, req.ID, categoryToUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (i *categoryImpl) DeleteCategory(ctx context.Context, id *int) error {
	err := i.CategoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
