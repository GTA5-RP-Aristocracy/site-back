package handler

import "github.com/GTA5-RP-Aristocracy/site-back/product"

type (
	// Service the product handler.
	Service interface {
		Create(product product.Product) error
		List(filter product.ProductFilter) ([]product.Product, error)
		Get(id int) (product.Product, error)
		Update(product product.Product) error
		Delete(id int) error
	}
)
