package user

// This file contains user repository related code.

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	// Repository represents the user repository.
	repository struct {
		db *sql.DB
	}
)

// NewRepository creates a new user repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

// Create inserts a new user into the repository.
func (r *repository) Create(user User) error {
	_, err := r.db.Exec("INSERT INTO user_storage (id,email, name, password) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.Name, user.Password)
	return err
}

// FindByEmail returns a user by email.
func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, email, name, password, created, updated FROM user_storage WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Created, &user.Updated)
	if errors.Is(err,sql.ErrNoRows){
		return User{},ErrNotFound
	}

	return user, err
}

// FindByID returns a user by id.
func (r *repository) FindByID(id uuid.UUID) (User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, email, name, password, created, updated FROM user_storage WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Created, &user.Updated)
	return user, err
}

// FindAll returns all users.
func (r *repository) FindAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, email, name, password, created, updated FROM user_storage")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Created, &user.Updated); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
