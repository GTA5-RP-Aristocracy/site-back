package repository

import (
	"database/sql"
	"fmt"

	"github.com/GTA5-RP-Aristocracy/site-back/product"
)

type (
	// Repository describes the persistence on product model.
	Repository struct {
		db *sql.DB
	}
)

// New creates a new product repository.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new product.
func (r *Repository) Create(product product.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name, price, description, picture) VALUES ($1, $2, $3, $4)", product.Name, product.Price, product.Description, product.Picture)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

// List returns a list of products.
func (r *Repository) List(filter product.ProductFilter) ([]product.Product, error) {
	q := "SELECT id, name, price, description, picture, created, updated FROM products"
	args := []interface{}{}
	if filter.Limit > 0 {
		q += " LIMIT $1"
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		q += " OFFSET $2"
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Picture, &p.Created, &p.Updated)
		products = append(products, p)
	}

	return products, nil
}

// Get returns a product by its ID.
func (r *Repository) Get(id int) (product.Product, error) {
	var p product.Product
	err := r.db.QueryRow("SELECT id, name, price, description, picture, created, updated FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Picture, &p.Created, &p.Updated)
	if err != nil {
		return product.Product{}, fmt.Errorf("get product: %w", err)
	}

	return p, nil
}

// Update updates a product.
func (r *Repository) Update(product product.Product) error {
	_, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, description = $3, picture = $4 WHERE id = $5",
		product.Name, product.Price, product.Description, product.Picture, product.ID,
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

// Delete deletes a product by its ID.
func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}
