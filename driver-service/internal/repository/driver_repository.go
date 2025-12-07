package repository

import (
	"context"
	"time"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// driverRepository interface defines database operations
type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) (string, error)
}

// driverRepositoryImpl implements DriverRepository for mongodb
type driverRepositoryImpl struct {
	collection *mongo.Collection
}

// NewDriverRepository creates a new repository instance
func NewDriverRepository(db *mongo.Database) DriverRepository {
	return &driverRepositoryImpl{
		collection: db.Collection("drivers"),
	}
}

// Create inserts a new driver into the database
func (r *driverRepositoryImpl) Create(ctx context.Context, driver *models.Driver) (string, error) {
	// set timestamps
	now := time.Now()
	driver.CreatedAt = now
	driver.UpdatedAt = now

	// insert into mongodb
	result, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return "", err
	}

	// return the generated object id as hex string
	oid, _ := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}
