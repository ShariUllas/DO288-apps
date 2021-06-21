package service

import (
	"catalogue/api/internal/core/data"
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/core/response"
	"context"
)

// Product is service interface
type Product interface {
	GetProduct(ctx context.Context) ([]*response.Product, error)
	AddProduct(ctx context.Context, category *request.Product) (response.AddProduct, error)
}

type productImpl struct {
	data.ProductRepo
}

// NewProductService is
func NewProductService(productRepo data.ProductRepo) Product {
	return &productImpl{ProductRepo: productRepo}
}

func (i *productImpl) AddProduct(ctx context.Context, category *request.Product) (response.AddProduct, error) {
	productID, err := i.ProductRepo.AddProduct(ctx, category)
	if err != nil {
		return response.AddProduct{}, err
	}
	return response.NewAddProduct(productID), nil
}

func (i *productImpl) GetProduct(ctx context.Context) ([]*response.Product, error) {
	products, err := i.ProductRepo.GetProduct(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewListProductResponse(products), nil
}
