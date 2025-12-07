package service

import (
	"context"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"github.com/eneszeyt/bitaksi-driver-service/internal/repository"
)

// DriverService defines business logic
type DriverService interface {
	CreateDriver(ctx context.Context, driver *models.Driver) (string, error)
	UpdateDriver(ctx context.Context, id string, driver *models.Driver) error
	ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, error)
}

type driverServiceImpl struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverServiceImpl{repo: repo}
}

func (s *driverServiceImpl) CreateDriver(ctx context.Context, driver *models.Driver) (string, error) {
	return s.repo.Create(ctx, driver)
}

func (s *driverServiceImpl) UpdateDriver(ctx context.Context, id string, driver *models.Driver) error {
	return s.repo.Update(ctx, id, driver)
}

// ListDrivers logic
func (s *driverServiceImpl) ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, error) {
	// default values logic could be here if needed
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return s.repo.List(ctx, page, pageSize)
}
