package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/entity"
)

// DigitalAssetCollectionRepository struct
type DigitalAssetCollectionRepository struct {
	dbMongo *mongo.Database
}

// NewDigitalAssetCollectionRepository create repository
func NewDigitalAssetCollectionRepository(dbMongo *mongo.Database) *DigitalAssetCollectionRepository {
	return &DigitalAssetCollectionRepository{dbMongo: dbMongo}
}

func (r *DigitalAssetCollectionRepository) GetCollectionByOrgID(orgID *string) (*entity.DigitalAssetCollection, error) {
	oID, err := primitive.ObjectIDFromHex(*orgID)
	if err != nil {
		return nil, err
	}

	var collection entity.DigitalAssetCollection
	err = r.dbMongo.Collection(collection.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "org_id", Value: oID}}).Decode(&collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (r *DigitalAssetCollectionRepository) GetCollectionByID(cID *string) (*entity.DigitalAssetCollection, error) {
	colID, err := primitive.ObjectIDFromHex(*cID)
	if err != nil {
		return nil, err
	}
	var dac entity.DigitalAssetCollection
	err = r.dbMongo.Collection(dac.CollectionName()).FindOne(context.TODO(), bson.M{"_id": colID}).Decode(&dac)
	if err != nil {
		return nil, err
	}

	return &dac, nil
}
