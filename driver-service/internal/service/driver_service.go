package service

import (
	"context"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"github.com/eneszeyt/bitaksi-driver-service/internal/repository"
)

// DriverService defines business logic
type DriverService interface {
	CreateDriver(ctx context.Context, driver *models.Driver) (string, error)
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
