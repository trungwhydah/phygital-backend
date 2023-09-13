package connections

import (
	"context"
	"fmt"
	"log"
	"time"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/common/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() (db *mongo.Database) {
	var (
		currentEntry     int
		ctx              = context.Background()
		maxEntry         = config.C.Mongo.MaxRetry
		err              error
		client           *mongo.Client
		connectionString = config.C.Mongo.MongoURLDbString
	)

	for {
		currentEntry++

		if currentEntry > maxEntry {
			log.Fatalf("MongoDB connection errors. Exit after try %d times", maxEntry)
		}

		client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
		if err != nil {
			logger.LogError(fmt.Sprintf("%s. %d times try", err.Error(), currentEntry))
			time.Sleep(time.Duration(config.C.Mongo.IntervalTime) * time.Second)

			continue
		}

		err = client.Connect(ctx)
		if err != nil {
			logger.LogError(fmt.Sprintf("Connect to DB fail: %s", err.Error()))
		}

		db = client.Database(config.C.Mongo.DatabaseName)

		return
	}
}
