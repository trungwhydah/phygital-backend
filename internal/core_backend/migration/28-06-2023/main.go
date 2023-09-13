package main

import (
	"context"
	"log"

	// "backend-service/internal/core_backend/migration/28-06-2023/organization"
	// "backend-service/internal/core_backend/migration/28-06-2023/product"
	// "backend-service/internal/core_backend/migration/28-06-2023/productItem"
	"backend-service/internal/core_backend/migration/28-06-2023/tag"
	// "backend-service/internal/core_backend/migration/28-06-2023/template_and_mapping"
	// "backend-service/internal/core_backend/migration/28-06-2023/user"
	// "backend-service/internal/core_backend/migration/28-06-2023/verification"
	// "backend-service/internal/core_backend/migration/28-06-2023/webPage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SOURCE_DB      = "phygital-staging"
	DESTINATION_DB = "phygital-staging-20230713"
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
	DestinationDB := client.Database(DESTINATION_DB)

	log.Println("Data migration starting...")

	//Comment if you don't want to migrate specific database
	tag.MigrateCollection(SourceDB, DestinationDB)
	// organization.MigrateCollection(SourceDB, DestinationDB)
	// productItem.MigrateCollection(SourceDB, DestinationDB)
	// webPage.MigrateCollection(SourceDB, DestinationDB)
	// verification.MigrateCollection(SourceDB, DestinationDB)
	// user.MigrateCollection(SourceDB, DestinationDB)
	// product.MigrateCollection(SourceDB, DestinationDB)
	// template_and_mapping.MigrateCollection(SourceDB, DestinationDB)

	log.Println("Data migration complete.")
}
