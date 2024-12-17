package user

import (
	"time"

	"github.com/google/uuid"
)

// This file defines the user model, e.g. the struct that represents a user.

const (
	// RoleAdmin represents an admin user.
	RoleAdmin Role = "admin"
	// RoleUser represents a regular user.
	RoleUser Role = "user"
	// RoleModerator represents a moderator user.
	RoleModerator Role = "moderator"
)

type (
	Role string
	// User represents a user account in the system.
	User struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Name     string    `json:"name"`
		Password string    `json:"password"`
		Role     Role      `json:"role"`
		Blocked  bool      `json:"blocked"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
	}

	UserFilter struct {
		Roles  []Role
		Limit  int
		Offset int
	}

	FieldsToUpdate struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     Role   `json:"role"`
		Blocked  bool   `json:"blocked"`
	}
)
