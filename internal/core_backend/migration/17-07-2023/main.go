package main

import (
	"context"
	"log"

	auto_sync_total_item "backend-service/internal/core_backend/migration/17-07-2023/auto-sync-total-item"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SOURCE_DB      = "phygital-staging-20230717"
	DESTINATION_DB = "phygital-staging-20230717-1"
	MONGO_URI      = ""
)

func main() {
	// Connect to the MongoDB instance
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Access the source database and collection
	SourceDB := client.Database(SOURCE_DB)

	// Access the destination database and collection
	// DestinationDB := client.Database(DESTINATION_DB)

	log.Println("Data migration starting...")

	//Comment if you don't want to migrate specific database
	// move_template_to_product.MigrateCollection(SourceDB, DestinationDB)
	// template.MigrateCollection(SourceDB, DestinationDB)
	auto_sync_total_item.SyncTotalItems(SourceDB)

	log.Println("Data migration complete.")
}
