package user

// This file contains the user service implementation.

type (
	// service implements the Service interface.
	service struct {
		repo Repository
	}
)

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &service{repo}
}

// Signup creates a new user account.
func (s *service) Signup(email, name, password string) error {
	// Check if the email is already registered.
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return ErrEmailExists
	}

	// Create a new user.
	user := User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	return s.repo.Create(user)
}

// Signin checks the email and password and returns a user.
func (s *service) Signin(email, password string) (User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return User{}, err
	}

	// TODO: Use a secure password hashing algorithm.
	if user.Password != password {
		return User{}, ErrNotFound
	}
	return user, nil
}

// Get fetches a user by id.
func (s *service) Get(id int) (User, error) {
	return s.repo.FindByID(id)
}

// List fetches all users.
func (s *service) List() ([]User, error) {
	return s.repo.FindAll()
}
