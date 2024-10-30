package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

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
	if err !=ErrNotFound{
		return fmt.Errorf("error get email:%w",err)
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
func (s *service) Get(id uuid.UUID) (User, error) {
	return s.repo.FindByID(id)
}

// List fetches all users.
func (s *service) List() ([]User, error) {
	return s.repo.FindAll()
}

// check passw and hash sum
func (s *service) checkPasswordHash(password, encodedHash string) (bool, error) {
	parts :=strings.Split(encodedHash,"$")
	if len(parts)!=2{
		return false, fmt.Errorf("invalid hash format")
	}
	passwdBase64 := parts[0]
	hashBase64 := parts[1]
	
	passwd, err := base64.RawStdEncoding.DecodeString(passwdBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(hashBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
	generatedHash := argon2.IDKey([]byte(password), passwd, 1, 64*1024, 4, 32)

	return string(expectedHash) == string(generatedHash), nil
}


// hashed password user
func (s *service) passHashed(password string) (string, error) {
	passwd := make([]byte, 16)
	_, err := rand.Read(passwd)
	if err != nil {
		return "", fmt.Errorf("error generating password: %w", err)
	}

	hash := argon2.IDKey([]byte(password), passwd, 1, 64*1024, 4, 32)
	passBase64 := base64.RawStdEncoding.EncodeToString(passwd)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", passBase64, hashBase64), nil
}
