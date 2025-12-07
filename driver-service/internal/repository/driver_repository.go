package repository

import (
	"context"
	"errors"
	"time"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DriverRepository defines database operations
type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) (string, error)
	Update(ctx context.Context, id string, driver *models.Driver) error
}

type driverRepositoryImpl struct {
	collection *mongo.Collection
}

func NewDriverRepository(db *mongo.Database) DriverRepository {
	return &driverRepositoryImpl{
		collection: db.Collection("drivers"),
	}
}

// Create inserts a new driver
func (r *driverRepositoryImpl) Create(ctx context.Context, driver *models.Driver) (string, error) {
	now := time.Now()
	driver.CreatedAt = now
	driver.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return "", err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

// Update modifies an existing driver
func (r *driverRepositoryImpl) Update(ctx context.Context, id string, driver *models.Driver) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	// update query: set updated_at and other fields
	driver.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"firstName": driver.FirstName,
			"lastName":  driver.LastName,
			"plate":     driver.Plate,
			"taxiType":  driver.TaxiType,
			"location":  driver.Location,
			"updatedAt": driver.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("driver not found")
	}

	return nil
}
