package user

import "github.com/google/uuid"

// This file defines the user related interfaces.

type (
	// Service represents the user service interface.
	Service interface {
		// Signup creates a new user account.
		Signup(email, name, password string) error
		// Signin checks the email and password and returns a user.
		Signin(email, password string) (User, error)
		// Get fetches a user by id.
		Get(id uuid.UUID) (User, error)
		// List fetches all users.
		List() ([]User, error)
	}

	// Repository represents the user repository interface.
	Repository interface {
		// Create inserts a new user into the repository.
		Create(user User) error
		// FindByEmail returns a user by email.
		FindByEmail(email string) (User, error)
		// FindByID returns a user by id.
		FindByID(id uuid.UUID) (User, error)
		// FindAll returns all users.
		FindAll() ([]User, error)
	}

	// Verifier represents the reCAPTCHA verifier interface.
	Verifier interface {
		// Verify verifies the reCAPTCHA response.
		Verify(response, remoteip string) (bool, error)
	}
)
