package user

// This file contains user repository related code.

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	err := r.db.QueryRow("SELECT id, email, name, password, created, updated, role, blocked FROM user_storage WHERE id = $1", id).
		Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Created, &user.Updated, &user.Role, &user.Blocked)
	return user, err
}

// FindAll returns all users.
func (r *repository) FindAll(filter UserFilter) ([]User, error) {
	q := "SELECT id, email, name, password, created, updated, blocked, role FROM user_storage"

	count := 0
	args := []interface{}{}
	if len(filter.Roles) > 0 {
		roles := make([]string, len(filter.Roles))
		for i, role := range filter.Roles {
			roles[i] = strconv.Itoa(int(role))
		}
		q += fmt.Sprintf(" WHERE role IN ('%s')", strings.Join(roles, "','"))
	}

	q += " ORDER BY created DESC"

	if filter.Limit > 0 {
		count++
		q += fmt.Sprintf(" LIMIT $%d", count)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		count++
		q += fmt.Sprintf(" OFFSET $%d", count)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Password, &user.Created, &user.Updated, &user.Blocked, &user.Role,
		); err != nil {
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

	count := 0
	if fields.Name != "" {
		count++
		q += fmt.Sprintf(", name = $%d", count)
		args = append(args, fields.Name)
	}

	if fields.Email != "" {
		count++
		q += fmt.Sprintf(", email = $%d", count)
		args = append(args, fields.Email)
	}

	if fields.Password != "" {
		count++
		q += fmt.Sprintf(", password = $%d", count)
		args = append(args, fields.Password)
	}

	if fields.Role != 0 {
		count++
		q += fmt.Sprintf(", role = $%d", count)
		args = append(args, fields.Role)
	}

	count++
	q += fmt.Sprintf(", blocked = $%d", count)
	args = append(args, fields.Blocked)

	count++
	q += fmt.Sprintf(" WHERE id = $%d", count)
	args = append(args, id)

	_, err := r.db.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
