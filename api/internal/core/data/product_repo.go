package data

import (
	"catalogue/api/internal/core/request"
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// ProductRepo is repo interface for product
type ProductRepo interface {
	GetProduct(ctx context.Context) ([]*Product, error)
	AddProduct(ctx context.Context, variant *request.Product) (*int, error)
}

type productRepoImpl struct {
	*sql.DB
}

// NewProductRepo is
func NewProductRepo(db *sql.DB) ProductRepo {
	return &productRepoImpl{DB: db}
}

// Product -
type Product struct {
	ID          int
	Name        string
	Category    string
	Description *string
	ImageURL    *string
}

const addProductQuery = `
INSERT INTO product
(
	name,
	category_id,
	description,
	image_url,
	created_at,
	updated_at
)
VALUES
(
	$1,$2,$3,$4,$5,$6
)
returning id;
`

func (i *productRepoImpl) AddProduct(ctx context.Context, product *request.Product) (*int, error) {
	var productID *int
	err := i.QueryRowContext(ctx, addProductQuery, product.Name, product.CategoryID, product.Description, product.ImageURL, time.Now(), time.Now()).Scan(&productID)
	if err != nil {
		return nil, errors.Wrapf(err, "adding product in to db failed")
	}
	return productID, nil
}

const getProductQuery = `
SELECT p.id,
       p.name,
	   p.description,
	   p.image_url,
       c.name AS category
FROM   product AS p
       LEFT JOIN category AS c
              ON p.category_id = c.id  
`

func (i *productRepoImpl) GetProduct(ctx context.Context) ([]*Product, error) {
	var products []*Product
	rows, err := i.QueryContext(ctx, getProductQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "get product query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImageURL, &product.Category)
		if err != nil {
			return nil, errors.Wrapf(err, "get product query scan failed")
		}
		products = append(products, &product)
	}
	return products, nil
}
