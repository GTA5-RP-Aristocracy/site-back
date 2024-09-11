package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/argon2"
)

// This file contains user related http handlers.

const (
	pathRoot   = "/user"
	pathList   = "/list"
	pathSignup = "/signup"
)

type (
	// Handler represents a set of http handlers for managing users.
	Handler struct {
		service Service
	}
)

// NewHandler creates a new user http handler.
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// RegisterUserRouter registers user routes.
func (h *Handler) RegisterUserRouter(externalRouter chi.Router) {
	r := chi.NewRouter()
	r.Post(pathSignup, h.Signup)
	r.Get(pathList, h.List)
	externalRouter.Mount(pathRoot, r)
	
}

// Signup handles user signup request.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	email := r.FormValue("email")
	name := r.FormValue("name")
	password := r.FormValue("password")

	// Create a new user.
	if err := h.service.Signup(email, name, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List handles user list request.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Fetch all users.
	users, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response.
	writeJSON(w, users)
}

// writeJSON writes the response in JSON format.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Write stas (hashed password user)
func (h *Handler) passHashed(password string) (string, error) {
	passwd := make([]byte,16)
	_, err := rand.Read(passwd)
	if err != nil{
		return "", fmt.Errorf("Error generating password: %w", err)
	}
	hash := argon2.IDKey([]byte(password), passwd,1,64*1024,4,32 )
	passBase64 := base64.RawStdEncoding.EncodeToString(passwd)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s$%s",passBase64,hashBase64), nil
	} 
	
	

	


// Write stas (check passw and hash sum)
func checkPasswordHash(password, encodedHash string)(bool, error){
	var passwdBase64, hashBase64 string
	_, err := fmt.Sscanf(encodedHash, "%s$%s",&passwdBase64,&hashBase64)
	if err != nil{
		return false,err
	}

	passwd, err := base64.RawStdEncoding.DecodeString(passwdBase64)
	if err != nil{
		return false,err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(hashBase64)
	if err != nil{
		return false, err
	}
	generatedHash := argon2.IDKey([]byte(password), passwd, 1, 64*1024, 4, 32)

	return string(expectedHash) == string(generatedHash), nil
}