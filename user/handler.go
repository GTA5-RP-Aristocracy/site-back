package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	siteback "github.com/GTA5-RP-Aristocracy/site-back"
	"github.com/GTA5-RP-Aristocracy/site-back/middlware"
	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// This file contains user related http handlers.

const (
	pathRoot    = "/user"
	pathList    = "/list"
	pathSignup  = "/signup"
	pathSignin  = "/signin"
	pathVerify  = "/verify"
	pathBlock   = "/block"
	pathUnblock = "/unblock"
	pathUpdate  = "/update"
)

type (
	// Handler represents a set of http handlers for managing users.
	Handler struct {
		service     Service
		verifier    Verifier
		tokenSecret string
	}

	VerifierRequest struct {
		Token string `json:"token"`
	}
)

// NewHandler creates a new user http handler.
func NewHandler(service Service, verifier Verifier, tokenSecret string) *Handler {
	return &Handler{
		service:     service,
		verifier:    verifier,
		tokenSecret: tokenSecret,
	}
}

// RegisterUserRouter registers user routes.
func (h *Handler) RegisterUserRouter(externalRouter chi.Router) {
	r := chi.NewRouter()
	r.Post(pathSignup, h.Signup)
	r.Post(pathSignin, h.Signin)
	r.Get(pathList, middlware.AuthMiddleware(h.tokenSecret, RoleMiddleware(RoleModerator, h.service, h.List)))
	r.Post(pathVerify, h.Verify)
	r.Get(pathBlock, middlware.AuthMiddleware(h.tokenSecret, RoleMiddleware(RoleModerator, h.service, h.Block)))
	r.Get(pathUnblock, middlware.AuthMiddleware(h.tokenSecret, RoleMiddleware(RoleModerator, h.service, h.Unblock)))
	r.Post(pathUpdate, middlware.AuthMiddleware(h.tokenSecret, RoleMiddleware(RoleUser, h.service, h.Update)))
	r.Get("/", middlware.AuthMiddleware(h.tokenSecret, RoleMiddleware(RoleUser, h.service, h.Get)))

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
	var (
		filter UserFilter
		err    error
	)

	// Parse the roles.
	roles := r.URL.Query()["role"]
	for _, role := range roles {
		r := RoleFromString(role)
		if !IsValidRole(r) {
			siteback.WriteError(w, fmt.Errorf("invalid role: %s", role), http.StatusBadRequest)
			return
		}

		filter.Roles = append(filter.Roles, r)
	}

	// Parse the limit.
	if limit := r.URL.Query().Get("limit"); limit != "" {
		filter.Limit, err = parseInt(limit, 10, 64)
		if err != nil {
			siteback.WriteError(w, err, http.StatusBadRequest)
			return
		}
	}

	// Parse the offset.
	if offset := r.URL.Query().Get("offset"); offset != "" {
		filter.Offset, err = parseInt(offset, 10, 64)
		if err != nil {
			siteback.WriteError(w, err, http.StatusBadRequest)
			return
		}
	}

	// Fetch all users.
	users, err := h.service.List(filter)
	if err != nil {
		siteback.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// Write the response.
	siteback.WriteJSON(w, users)
}

// parseInt parses the string s into an integer type.
func parseInt(s string, base int, bitSize int) (int, error) {
	n, err := strconv.ParseInt(s, base, bitSize)
	return int(n), err
}

// TODO move to separate file
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserResponse struct {
	ID    uuid.UUID
	Email string
	Name  string
	Token string
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, token, err := h.service.Signin(email, password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Token: token,
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
		if err == ErrNotFound {
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
	siteback.WriteJSON(w, map[string]bool{"success": ok})
}

// Block handles user block request.
func (h *Handler) Block(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}

	// Parse the UUID.
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Block the user.
	if err := h.service.Block(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Unblock handles user unblock request.
func (h *Handler) Unblock(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}

	// Parse the UUID.
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Unblock the user.
	if err := h.service.Unblock(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Update handles user update request.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "User UUID is required", http.StatusBadRequest)
		return
	}

	// Parse the UUID.
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Parse the request.
	var fields FieldsToUpdate
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the user.
	if err := h.service.Update(id, fields); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RoleMiddleware(role Role, userService Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idFromCtx := siteback.UserIDFromContext(r.Context())
		if idFromCtx == uuid.Nil {
			siteback.WriteError(w, errors.New("missing user id"), http.StatusUnauthorized)
			return
		}

		u, err := userService.Get(idFromCtx)
		if err != nil {
			siteback.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		if u.Role > role {
			siteback.WriteError(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
