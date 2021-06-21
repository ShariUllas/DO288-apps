package service

import (
	"catalogue/api/internal/core/data"
	"catalogue/api/internal/core/request"
	"catalogue/api/internal/core/response"
	"context"
	"fmt"
)

// Variant is service interface
type Variant interface {
	GetVariant(ctx context.Context) ([]*response.Variant, error)
	AddVariant(ctx context.Context, variant *request.Variant) (*int, error)
	EditVariant(ctx context.Context, variant *request.Variant) error
	DeleteVariant(ctx context.Context, id *int) error
}

type variantImpl struct {
	data.VariantRepo
}

var (
	// ErrVariantNotFound is returned whne variant is not found
	ErrVariantNotFound = fmt.Errorf("variant not found")
)

// NewVariantService - service for variant related operations
func NewVariantService(variantRepo data.VariantRepo) Variant {
	return &variantImpl{VariantRepo: variantRepo}
}

func (i *variantImpl) GetVariant(ctx context.Context) ([]*response.Variant, error) {
	variants, err := i.VariantRepo.GetVariant(ctx)
	if err != nil {
		return nil, err
	}
	return response.NewVariantResponse(variants), nil
}

func (i *variantImpl) AddVariant(ctx context.Context, variant *request.Variant) (*int, error) {
	variantID, err := i.VariantRepo.AddVariant(ctx, variant)
	if err != nil {
		return nil, err
	}
	return variantID, nil
}

func (i *variantImpl) EditVariant(ctx context.Context, req *request.Variant) error {
	// request may contain partial data so we overwrite the existing with it
	variantToUpdate, err := i.VariantRepo.GetVariantByID(ctx, req.ID)
	if err != nil {
		if err == data.ErrVariantNotFound {
			return ErrVariantNotFound
		}
		return err
	}
	if len(req.Color) == 0 {
		variantToUpdate.Color = &req.Color
	}
	if req.DiscountPrice <= 0 {
		variantToUpdate.DiscountPrice = &req.DiscountPrice
	}
	if req.MRP <= 0 {
		variantToUpdate.MRP = req.MRP
	}
	if len(req.Name) != 0 {
		variantToUpdate.Name = &req.Name
	}
	if req.ProductID <= 0 {
		variantToUpdate.ProductID = req.ProductID
	}
	if req.Size <= 0 {
		variantToUpdate.Size = &req.Size
	}

	err = i.VariantRepo.Edit(ctx, req.ID, variantToUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (i *variantImpl) DeleteVariant(ctx context.Context, id *int) error {
	err := i.VariantRepo.DeleteVariant(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
