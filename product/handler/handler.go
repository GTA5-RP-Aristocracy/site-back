package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GTA5-RP-Aristocracy/site-back/product"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

const (
	pathRoot   = "/product"
	pathCreate = "/"
	pathList   = "/"
	pathGet    = "/{id}"
	pathUpdate = "/{id}"
	pathDelete = "/{id}"
)

type (
	// Handler is the product handler.
	Handler struct {
		service Service
		log     zerolog.Logger
	}
)

// RegisterProductRouter registers product routes.
func (h *Handler) RegisterProductRouter(externalRouter chi.Router) {
	r := chi.NewRouter()
	r.Post(pathCreate, h.Create)
	r.Get(pathList, h.List)
	r.Get(pathGet, h.Get)
	r.Put(pathUpdate, h.Update)
	r.Delete(pathDelete, h.Delete)

	externalRouter.Mount(pathRoot, r)
}

// NewHandler creates a new product handler.
func NewHandler(service Service, log zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// Create handles product creation request.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	var product product.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.log.Err(err).Msg("failed to parse product")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Create the product.
	if err := h.service.Create(product); err != nil {
		h.log.Err(err).Msg("failed to create product")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Respond.
	w.WriteHeader(http.StatusCreated)
}

// List handles product list request.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	var filter product.ProductFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		h.log.Err(err).Msg("failed to parse product filter")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Fetch all products.
	products, err := h.service.List(filter)
	if err != nil {
		h.log.Err(err).Msg("failed to fetch products")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Respond.
	if err := json.NewEncoder(w).Encode(products); err != nil {
		h.log.Err(err).Msg("failed to encode products")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Get handles product get request.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Err(err).Msg("failed to parse product ID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Fetch the product.
	product, err := h.service.Get(id)
	if err != nil {
		h.log.Err(err).Msg("failed to fetch product")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Respond.
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.log.Err(err).Msg("failed to encode product")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Update handles product update request.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Err(err).Msg("failed to parse product ID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var product product.Product
	product.ID = id
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.log.Err(err).Msg("failed to parse product")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Update the product.
	if err := h.service.Update(product); err != nil {
		h.log.Err(err).Msg("failed to update product")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

// Delete handles product delete request.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse the request.
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Err(err).Msg("failed to parse product ID")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Delete the product.
	if err := h.service.Delete(id); err != nil {
		h.log.Err(err).Msg("failed to delete product")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}
