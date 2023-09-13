package repository

import (
	"context"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// VerificationRepository struct
type VerificationRepository struct {
	dbMongo *mongo.Database
}

// NewVerificationRepository create repository
func NewVerificationRepository(dbMongo *mongo.Database) *VerificationRepository {
	return &VerificationRepository{dbMongo: dbMongo}
}

func (r *VerificationRepository) SaveVerifition(ver *entity.Verification) (*entity.Verification, error) {
	result, err := r.dbMongo.Collection(ver.CollectionName()).InsertOne(context.TODO(), ver)
	if err != nil {
		return nil, err
	}
	ver.ID = result.InsertedID.(primitive.ObjectID)

	return ver, nil
}
