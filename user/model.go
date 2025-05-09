package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// This file defines the user model, e.g. the struct that represents a user.

const (
	RoleAdmin Role = iota + 1
	RoleModerator
	RoleUser

	RoleAdminStr     = "admin"
	RoleModeratorStr = "moderator"
	RoleUserStr      = "user"
)

type (
	Role int
	// User represents a user account in the system.
	User struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Name     string    `json:"name"`
		Password string    `json:"password,omitempty"`
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

func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleModerator:
		return true
	}
	return false
}

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return RoleAdminStr
	case RoleModerator:
		return RoleModeratorStr
	case RoleUser:
		return RoleUserStr
	}
	return ""
}

func RoleFromString(role string) Role {
	switch role {
	case RoleAdminStr:
		return RoleAdmin
	case RoleModeratorStr:
		return RoleModerator
	case RoleUserStr:
		return RoleUser
	}
	return 1000
}

func (r Role) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var role string
	if err := json.Unmarshal(data, &role); err != nil {
		return fmt.Errorf("invalid role: %w", err)
	}

	*r = RoleFromString(role)
	return nil
}
