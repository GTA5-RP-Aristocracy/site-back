package user

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// This file contains user related http handlers.

const (
	pathRoot   = "/user"
	pathList   = "/list"
	pathSignup = "/signup"
	pathSignin = "/signin"
	pathVerify = "/verify"
)

type (
	// Handler represents a set of http handlers for managing users.
	Handler struct {
		service  Service
		verifier Verifier
	}

	VerifierRequest struct {
		Token string `json:"token"`
	}
)

// NewHandler creates a new user http handler.
func NewHandler(service Service, verifier Verifier) *Handler {
	return &Handler{
		service:  service,
		verifier: verifier,
	}
}

// RegisterUserRouter registers user routes.
func (h *Handler) RegisterUserRouter(externalRouter chi.Router) {
	r := chi.NewRouter()
	r.Post(pathSignup, h.Signup)
	r.Post(pathSignin, h.Signin)
	r.Get(pathList, h.List)
	r.Post(pathVerify, h.Verify)

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

// TODO move to separate file
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserResponse struct {
	Email string
	Name  string
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.Signin(email, password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := UserResponse{
		Email: user.Email,
		Name:  user.Name,
	}

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: user.ID.String(),
	})

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

// TODO move to separate file
var ErrUserNotFound = errors.New("user not found")

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")

	if uuidStr == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}

	parsUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	user, err := h.service.Get(parsUUID)

	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to serialize user data", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	var req VerifierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the remote IP address.
	remoteip := r.RemoteAddr

	// Verify the reCAPTCHA response.
	ok, err := h.verifier.Verify(req.Token, remoteip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response.
	writeJSON(w, map[string]bool{"success": ok})
}
