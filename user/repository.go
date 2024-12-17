package user

// This file contains user repository related code.

import (
	"database/sql"
	"errors"
	"fmt"

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
func NewRepository(db *sql.DB) *repository {
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
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
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
func (r *repository) FindAll(filter UserFilter) ([]User, error) {
	q := "SELECT id, email, name, password, created, updated FROM user_storage"

	if len(filter.Roles) > 0 {
		q += " WHERE role = ANY($1)"
	}

	if filter.Limit > 0 {
		q += " LIMIT $1"
	}

	if filter.Offset > 0 {
		q += " OFFSET $2"
	}

	rows, err := r.db.Query(q)
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

// Update updates a user.
func (r *repository) Update(id uuid.UUID, fields FieldsToUpdate) error {
	q := "UPDATE user_storage SET updated = now()"
	args := []interface{}{}

	if fields.Name != "" {
		q += ", name = $1"
		args = append(args, fields.Name)
	}

	if fields.Email != "" {
		q += ", email = $2"
		args = append(args, fields.Email)
	}

	if fields.Password != "" {
		q += ", password = $3"
		args = append(args, fields.Password)
	}

	if fields.Role != "" {
		q += ", role = $4"
		args = append(args, fields.Role)
	}

	q += ", blocked = $5"
	args = append(args, fields.Blocked)

	q += " WHERE id = $6"
	args = append(args, id)

	_, err := r.db.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
