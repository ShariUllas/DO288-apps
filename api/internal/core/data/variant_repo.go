package data

import (
	"catalogue/api/internal/core/request"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// VariantRepo is repo interface for category
type VariantRepo interface {
	GetVariant(ctx context.Context) ([]*Variant, error)
	AddVariant(ctx context.Context, variant *request.Variant) (*int, error)
	GetVariantByID(ctx context.Context, variantID int) (*Variant, error)
	Edit(ctx context.Context, variantID int, editedVariant *Variant) error
	DeleteVariant(ctx context.Context, variantID *int) error
}

type variantRepoImpl struct {
	*sql.DB
}

// NewVariantRepo is
func NewVariantRepo(db *sql.DB) VariantRepo {
	return &variantRepoImpl{DB: db}
}

var (
	// ErrVariantNotFound is returned whne variant is not found
	ErrVariantNotFound = fmt.Errorf("variant not found")
)

// Variant - model for variant
type Variant struct {
	ID            int
	Name          *string
	MRP           int
	DiscountPrice *int
	Size          *int
	Color         *string
	Product       string
	ProductID     int
}

const getVariantQuery = `
SELECT v.id,
       v.NAME,
       v.mrp,
       v.discount_price,
       v.size,
       v.color,
       p.NAME AS product
FROM   variant v
       LEFT JOIN product AS p
              ON v.product_id = p.id
`

func (i *variantRepoImpl) GetVariant(ctx context.Context) ([]*Variant, error) {
	var variants []*Variant
	rows, err := i.QueryContext(ctx, getVariantQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "get variant query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var variant Variant
		err := rows.Scan(&variant.ID, &variant.Name, &variant.MRP, &variant.DiscountPrice, &variant.Size, &variant.Color, &variant.Product)
		if err != nil {
			return nil, errors.Wrapf(err, "get variant query scan failed")
		}
		variants = append(variants, &variant)
	}
	return variants, nil
}

const addVariantQuery = `
INSERT INTO variant
(
	name,
	product_id,
	mrp,
	discount_price,
	size,
	color,
	created_at,
	updated_at
)
VALUES
(
	$1,$2,$3,$4,$5,$6,$7,$8
)
returning id;
`

func (i *variantRepoImpl) AddVariant(ctx context.Context, variant *request.Variant) (*int, error) {
	var variantID *int
	err := i.QueryRowContext(ctx, addVariantQuery, variant.Name, variant.ProductID, variant.MRP, variant.DiscountPrice, variant.Size, variant.Color, time.Now(), time.Now()).Scan(&variantID)
	if err != nil {
		return nil, errors.Wrapf(err, "adding variant to db failed")
	}
	return variantID, nil
}

const getVariantByIDQuery = `
SELECT NAME,
       product_id,
       mrp,
       discount_price,
       size,
       color
FROM   variant
WHERE  id = $1
`

// GetVariantByID fetches a variant based on id
func (i *variantRepoImpl) GetVariantByID(ctx context.Context, variantID int) (*Variant, error) {
	var variant Variant
	if err := i.QueryRowContext(ctx, getVariantByIDQuery, variantID).Scan(
		&variant.Name,
		&variant.ProductID,
		&variant.MRP,
		&variant.DiscountPrice,
		&variant.Size,
		&variant.Color,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrVariantNotFound
		}
		return nil, errors.Wrap(err, "get variant query scan failed")
	}

	return &variant, nil
}

const queryVariantEdit = `
UPDATE variant
SET    NAME = $2,
       product_id = $3,
       mrp = $4,
       discount_price = $5,
       size = $6,
       color = $7,
       updated_at = $8
WHERE  id = $1
`

func (i *variantRepoImpl) Edit(ctx context.Context, variantID int, editedVariant *Variant) error {
	_, err := i.ExecContext(ctx, queryVariantEdit, variantID, editedVariant.Name, editedVariant.ProductID, editedVariant.MRP, editedVariant.DiscountPrice, editedVariant.Size, editedVariant.Color, time.Now())
	if err != nil {
		return errors.Wrap(err, "edit variant query failed")
	}
	return nil
}

const queryVariantDelete = `
UPDATE	variant
	SET	
		is_deleted = true,
		updated_at = $2
WHERE	id = $1;
`

// Delete soft-deletes an existing variant
func (i *variantRepoImpl) DeleteVariant(ctx context.Context, userID *int) error {
	_, err := i.ExecContext(ctx, queryVariantDelete, userID, time.Now())
	if err != nil {
		return errors.Wrap(err, "delete variant query failed")
	}

	return nil
}
