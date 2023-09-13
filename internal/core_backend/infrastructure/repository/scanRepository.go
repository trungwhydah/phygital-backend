package repository

import (
	"context"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ScanRepository struct
type ScanRepository struct {
	dbMongo *mongo.Database
}

// NewScanRepository create repository
func NewScanRepository(dbMongo *mongo.Database) *ScanRepository {
	return &ScanRepository{dbMongo: dbMongo}
}

// GetTagWithID
func (r *ScanRepository) GetTagWithID(chipID *string) (*entity.Tag, error) {
	var chip entity.Tag
	err := r.dbMongo.Collection(chip.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "tag_id", Value: &chipID}}).Decode(&chip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &chip, nil
}
