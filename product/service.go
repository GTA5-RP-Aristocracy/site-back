package product

type (
	// Service the product service.
	Service struct {
		repo Repository
	}
)

// NewService creates a new product service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new product.
func (s *Service) Create(product Product) error {
	return s.repo.Create(product)
}

// List returns a list of products.
func (s *Service) List(filter ProductFilter) ([]Product, error) {
	return s.repo.List(filter)
}

// Get returns a product by its ID.
func (s *Service) Get(id int) (Product, error) {
	return s.repo.Get(id)
}

// Update updates a product.
func (s *Service) Update(product Product) error {
	return s.repo.Update(product)
}

// Delete deletes a product by its ID.
func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
