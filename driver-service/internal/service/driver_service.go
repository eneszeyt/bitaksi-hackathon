package service

import (
	"context"
	"sort"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"github.com/eneszeyt/bitaksi-driver-service/internal/repository"
	"github.com/eneszeyt/bitaksi-driver-service/internal/utils"
)

// DriverService defines business logic
type DriverService interface {
	CreateDriver(ctx context.Context, driver *models.Driver) (string, error)
	UpdateDriver(ctx context.Context, id string, driver *models.Driver) error
	ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, error)
	FindNearby(ctx context.Context, lat, lon float64, taxiType string) ([]map[string]interface{}, error)
}

type driverServiceImpl struct {
	repo repository.DriverRepository
}

// NewDriverService creates service instance
func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverServiceImpl{repo: repo}
}

// CreateDriver implements the business logic for creating a driver
func (s *driverServiceImpl) CreateDriver(ctx context.Context, driver *models.Driver) (string, error) {
	return s.repo.Create(ctx, driver)
}

// UpdateDriver logic
func (s *driverServiceImpl) UpdateDriver(ctx context.Context, id string, driver *models.Driver) error {
	return s.repo.Update(ctx, id, driver)
}

// ListDrivers logic
func (s *driverServiceImpl) ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return s.repo.List(ctx, page, pageSize)
}

// FindNearby logic: filter by radius and sort by distance
func (s *driverServiceImpl) FindNearby(ctx context.Context, lat, lon float64, taxiType string) ([]map[string]interface{}, error) {
	// 1. Get candidate drivers from DB
	drivers, err := s.repo.Search(ctx, taxiType)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	// 2. Filter by distance (Haversine)
	for _, d := range drivers {
		dist := utils.CalculateDistance(lat, lon, d.Location.Lat, d.Location.Lon)

		// Radius check: 6 km
		if dist <= 6.0 {
			// Create a response object with distance
			res := map[string]interface{}{
				"id":         d.ID,
				"firstName":  d.FirstName,
				"lastName":   d.LastName,
				"plate":      d.Plate,
				"taxiType":   d.TaxiType,
				"location":   d.Location,
				"distanceKm": dist,
			}
			results = append(results, res)
		}
	}

	// 3. Sort by distance (nearest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i]["distanceKm"].(float64) < results[j]["distanceKm"].(float64)
	})

	return results, nil
}
