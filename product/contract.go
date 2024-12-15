package product

type (
	// Repository describes the persistence on product model.
	Repository interface {
		Create(product Product) error
		List(filter ProductFilter) ([]Product, error)
		Get(id int) (Product, error)
		Update(product Product) error
		Delete(id int) error
	}

	// ProductFilter represents the filtering options for listing products.
	ProductFilter struct {
		Limit  int
		Offset int
	}
)
