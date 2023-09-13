package product

import (
	"context"

	"backend-service/internal/marketplace/entity/product"
	"backend-service/pkg/common/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoInterface interface {
	FindOneByID(ctx context.Context, id *string) (*product.Product, error)
}

type MongoRepo struct {
	db       *mongo.Database
	collName string
}

func (r *MongoRepo) FindOneByID(ctx context.Context, id *string) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		logger.Errorw(
			"decode string to objectID err",
			"id", *id,
			"err", err,
		)

		return nil, err
	}

	var res product.Product
	err = r.db.Collection(r.collName).FindOne(ctx, bson.M{"_id": objectID}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		logger.Errorw(
			"find one by id fail",
			"id", *id,
			"err", err,
		)

		return nil, err
	}

	return &res, nil
}

func NewMongoRepo(
	db *mongo.Database,
) RepoInterface {
	return &MongoRepo{
		db:       db,
		collName: "products",
	}
}
