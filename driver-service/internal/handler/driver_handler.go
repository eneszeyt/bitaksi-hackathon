package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"github.com/eneszeyt/bitaksi-driver-service/internal/service"
)

type DriverHandler struct {
	service service.DriverService
}

func NewDriverHandler(service service.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

// CreateDriver handles POST requests
func (h *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var driver models.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateDriver(r.Context(), &driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// UpdateDriver handles PUT requests (e.g., /drivers/{id})
func (h *DriverHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// extract id from url path (simple parsing)
	// path: /drivers/12345 -> id: 12345
	id := strings.TrimPrefix(r.URL.Path, "/drivers/")
	if id == "" {
		http.Error(w, "missing driver id", http.StatusBadRequest)
		return
	}

	var driver models.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateDriver(r.Context(), id, &driver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}`))
}
