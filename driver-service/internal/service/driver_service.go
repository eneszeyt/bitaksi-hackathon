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

// UpdateDriver logic
func (s *driverServiceImpl) UpdateDriver(ctx context.Context, id string, driver *models.Driver) error {
	return s.repo.Update(ctx, id, driver)
}
