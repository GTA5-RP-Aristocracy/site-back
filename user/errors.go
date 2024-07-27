package user

// This file contains user related errors.

import "errors"

// Define custom errors.
var (
	ErrEmailExists = errors.New("user: email already exists")
	ErrNotFound    = errors.New("user: not found")
)
