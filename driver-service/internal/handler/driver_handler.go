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

// DriversRoot handles /drivers endpoint
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

// DriverByID handles /drivers/{id} endpoint
func (h *DriverHandler) DriverByID(w http.ResponseWriter, r *http.Request) {
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

// SearchNearby godoc
// @Summary      Find nearby drivers
// @Description  Calculates distance using Haversine formula and returns drivers within 6km radius
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        lat       query     number  true  "Latitude"
// @Param        lon       query     number  true  "Longitude"
// @Param        taxiType  query     string  false "Taxi Type (e.g. yellow, black)"
// @Success      200       {array}   models.Driver
// @Router       /drivers/nearby [get]
func (h *DriverHandler) SearchNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	taxiType := r.URL.Query().Get("taxiType")

	if latStr == "" || lonStr == "" {
		http.Error(w, "missing lat or lon parameters", http.StatusBadRequest)
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		http.Error(w, "invalid coordinates", http.StatusBadRequest)
		return
	}

	results, err := h.service.FindNearby(r.Context(), lat, lon, taxiType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// --- Private Helper Methods (Annotated for Swagger) ---

// createDriver godoc
// @Summary      Create a new driver
// @Description  Registers a new taxi driver in the database
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        driver  body      models.Driver  true  "Driver Information"
// @Success      201     {object}  map[string]string
// @Router       /drivers [post]
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

// updateDriver godoc
// @Summary      Update a driver
// @Description  Updates existing driver information by ID
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        id      path      string         true  "Driver ID"
// @Param        driver  body      models.Driver  true  "Driver Data"
// @Success      200     {object}  map[string]string
// @Router       /drivers/{id} [put]
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

// listDrivers godoc
// @Summary      List drivers
// @Description  Get all drivers with pagination
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number"
// @Param        pageSize  query     int  false  "Page size"
// @Success      200       {array}   models.Driver
// @Router       /drivers [get]
func (h *DriverHandler) listDrivers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	drivers, err := h.service.ListDrivers(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if drivers == nil {
		drivers = []models.Driver{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}
