package user

import (
	"time"

	"github.com/google/uuid"
)

// This file defines the user model, e.g. the struct that represents a user.

type (
	// User represents a user account in the system.
	User struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Name     string    `json:"name"`
		Password string    `json:"password"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
		
	}
)
