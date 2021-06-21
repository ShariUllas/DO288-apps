package data

import (
	"catalogue/api/internal/core/request"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// CategoryRepo is repo interface for category
type CategoryRepo interface {
	GetCategory(ctx context.Context) ([]*Category, error)
	AddCategory(ctx context.Context, category *request.Category) (*int, error)
	GetCategoryID(ctx context.Context, categoryID int) (*Category, error)
	EditCategory(ctx context.Context, categoryID int, editedCategory *Category) error
	DeleteCategory(ctx context.Context, categoryID *int) error
	ListCategory(ctx context.Context, id *int) ([]*Category, error)
}

type categoryRepoImpl struct {
	*sql.DB
}

// NewCategoryRepo is
func NewCategoryRepo(db *sql.DB) CategoryRepo {
	return &categoryRepoImpl{DB: db}
}

var (
	// ErrCategoryNotFound is returned whne variant is not found
	ErrCategoryNotFound = fmt.Errorf("category not found")
)

// Category -
type Category struct {
	ID       int
	Name     string
	Category *string
	ParentID int
}

const listCategoryQuery = `
SELECT id,
	   name
FROM category
WHERE parent_id =$1`

func (i *categoryRepoImpl) ListCategory(ctx context.Context, id *int) ([]*Category, error) {

	var categories []*Category
	rows, err := i.DB.QueryContext(ctx, listCategoryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

const addCategoryQuery = `
INSERT INTO category
(
	name,
	parent_id,
	created_at,
	updated_at
)
VALUES
(
	$1,$2,$3,$4
)
returning id;
`

func (i *categoryRepoImpl) AddCategory(ctx context.Context, category *request.Category) (*int, error) {
	var categoryID *int
	err := i.QueryRowContext(ctx, addCategoryQuery, category.Name, category.ParentID, time.Now(), time.Now()).Scan(&categoryID)
	if err != nil {
		return nil, errors.Wrapf(err, "adding category to db failed")
	}
	return categoryID, nil
}

const getCategoryQuery = `
SELECT c.id,
       c.name,
       p.name AS category
FROM   category AS c
       LEFT JOIN category AS p
              ON c.parent_id = p.id  
`

func (i *categoryRepoImpl) GetCategory(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	rows, err := i.QueryContext(ctx, getCategoryQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "get category query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Category)
		if err != nil {
			return nil, errors.Wrapf(err, "get category query scan failed")
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

const getCategoryByIDQuery = `
SELECT id,
       name,
	   parent_id
FROM   category
WHERE  id = $1
`

// GetCategoryID fetches a variant based on id
func (i *categoryRepoImpl) GetCategoryID(ctx context.Context, categoryID int) (*Category, error) {
	var category Category
	if err := i.QueryRowContext(ctx, getCategoryByIDQuery, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}
		return nil, errors.Wrap(err, "get category by id query scan failed")
	}

	return &category, nil
}

const editCategoryQuery = `
UPDATE category
SET    NAME = $2,
	   parent_id = $3,     
       updated_at = $8
WHERE  id = $1
`

func (i *categoryRepoImpl) EditCategory(ctx context.Context, categoryID int, editedCategory *Category) error {
	_, err := i.ExecContext(ctx, editCategoryQuery, categoryID, editedCategory.Name, editedCategory.ParentID, time.Now())
	if err != nil {
		return errors.Wrap(err, "edit category query failed")
	}
	return nil
}

const queryCategoryDelete = `
UPDATE	category
	SET	
		is_deleted = true,
		updated_at = $2
WHERE	id = $1;
`

// Delete soft-deletes an existing category
func (i *categoryRepoImpl) DeleteCategory(ctx context.Context, categoryID *int) error {
	_, err := i.ExecContext(ctx, queryCategoryDelete, categoryID, time.Now())
	if err != nil {
		return errors.Wrap(err, "delete category query failed")
	}

	return nil
}
