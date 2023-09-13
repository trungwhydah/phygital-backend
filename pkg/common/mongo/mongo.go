// Package mongo implements MongoDB connection.
package mongo

import (
	"context"

	config "backend-service/config/common"
	cmentity "backend-service/internal/common/entity"
	"backend-service/pkg/common/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/mgocompat"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
)

// NewMongoDatabase returns mongo database.
//
// cfg specific for service config.
// lc fx lifecycle to handle disconnect database when service down.
func New(cfg *config.Config, lc fx.Lifecycle) (*mongo.Database, error) {
	ctx := context.Background()

	client, err := mongo.Connect(
		ctx,
		options.Client().SetRegistry(mgocompat.Registry),
		options.Client().ApplyURI(cfg.Mongo.ConnURI),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	logger.Info("Connect to MongoDB successfully")

	db := client.Database(cfg.Mongo.DBName)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing MongoDB connection")
				return client.Disconnect(ctx)
			},
		},
	)

	return db, nil
}

func MapToBSON(val map[string]any) bson.D {
	newFilter := bson.D{}
	for k, v := range val {
		newFilter = append(newFilter, bson.E{Key: k, Value: v})
	}

	return newFilter
}

func ConvertFilterToBSON(filter map[string]any) bson.D {
	newFilter := bson.D{}
	for k, v := range filter {
		newFilter = append(newFilter, bson.E{Key: k, Value: v})
	}

	return newFilter
}

func NewID() cmentity.ID {
	return cmentity.ID(primitive.NewObjectID().Hex())
}

type countResult struct {
	Total int64 `json:"total" bson:"total"`
}

func CountDocuments(
	ctx context.Context,
	coll *mongo.Collection,
	pipeline mongo.Pipeline,
) int64 {
	countPipeline := pipeline
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	cursors, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		logger.Errorw(
			"cannot count document for",
			"collection", coll.Name(),
			"err", err,
		)

		return 0
	}

	var result []countResult
	if err := cursors.All(ctx, &result); err != nil {
		logger.Errorw(
			"decode when count document err",
			"collection", coll.Name(),
			"err", err,
		)

		return 0
	}

	if len(result) > 0 {
		return result[0].Total
	}

	return 0
}
