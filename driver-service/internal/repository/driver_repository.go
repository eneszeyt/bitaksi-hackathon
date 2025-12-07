package repository

import (
	"context"
	"errors"
	"time"

	"github.com/eneszeyt/bitaksi-driver-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DriverRepository defines database operations
type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) (string, error)
	Update(ctx context.Context, id string, driver *models.Driver) error
	List(ctx context.Context, page, pageSize int) ([]models.Driver, error)
	// new method :
	Search(ctx context.Context, taxiType string) ([]models.Driver, error)
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

// List returns a paginated list of drivers
func (r *driverRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]models.Driver, error) {
	// calculate skip count (e.g. page 1 -> skip 0, page 2 -> skip 20)
	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}) // newest first

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []models.Driver
	if err := cursor.All(ctx, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}

// Search returns drivers matching a criteria (e.g. taxi type)
func (r *driverRepositoryImpl) Search(ctx context.Context, taxiType string) ([]models.Driver, error) {
	filter := bson.M{}

	// if taxiType is provided, filter by it
	if taxiType != "" {
		filter["taxiType"] = taxiType
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []models.Driver
	if err := cursor.All(ctx, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}
