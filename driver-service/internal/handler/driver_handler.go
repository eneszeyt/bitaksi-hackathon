package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// DriversRoot handles /drivers endpoint (GET for List, POST for Create)
func (h *DriverHandler) DriversRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createDriver(w, r)
	case http.MethodGet:
		h.listDrivers(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// DriverByID handles /drivers/{id} endpoint (PUT for Update)
func (h *DriverHandler) DriverByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := strings.TrimPrefix(r.URL.Path, "/drivers/")
	if id == "" {
		http.Error(w, "missing driver id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updateDriver(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// private method: createDriver
func (h *DriverHandler) createDriver(w http.ResponseWriter, r *http.Request) {
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

// private method: updateDriver
func (h *DriverHandler) updateDriver(w http.ResponseWriter, r *http.Request, id string) {
	var driver models.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateDriver(r.Context(), id, &driver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// private method: listDrivers
func (h *DriverHandler) listDrivers(w http.ResponseWriter, r *http.Request) {
	// parse query params ?page=1&pageSize=20
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// service will handle defaults if zero
	drivers, err := h.service.ListDrivers(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ensure empty slice instead of null in json
	if drivers == nil {
		drivers = []models.Driver{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}
