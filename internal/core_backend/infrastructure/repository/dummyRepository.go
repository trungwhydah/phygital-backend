package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/entity"
)

// DummyRepository struct
type DummyRepository struct {
	dbMongo *mongo.Database
}

// NewDummyRepository create repository
func NewDummyRepository(dbMongo *mongo.Database) *DummyRepository {
	return &DummyRepository{dbMongo: dbMongo}
}

func (r *DummyRepository) GetDummyByID(DummyID *int) (*entity.Dummy, error) {
	var dummy entity.Dummy
	err := r.dbMongo.Collection(dummy.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "id", Value: &DummyID}}).Decode(&dummy)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}

		return nil, nil
	}

	return &dummy, nil
}
