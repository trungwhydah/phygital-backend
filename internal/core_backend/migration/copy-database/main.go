package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SOURCE_DB      = "phygital-staging-20230717-1"
	DESTINATION_DB = "phygital-staging-20230717-1"
	MONGO_URI      = ""
)

func main() {
	// Connect to the source MongoDB server
	sourceClient, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()
	err = sourceClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceClient.Disconnect(ctx)

	// Connect to the destination MongoDB server
	destinationClient, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()
	err = destinationClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer destinationClient.Disconnect(ctx)

	// Get the source and destination databases
	sourceDB := sourceClient.Database(SOURCE_DB)
	destinationDB := destinationClient.Database(DESTINATION_DB)

	// List the collections in the source database
	collections, err := sourceDB.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Copy each collection from the source to the destination database
	for _, collName := range collections {
		log.Println("Migrating", collName, "...")
		coll := sourceDB.Collection(collName)

		// Read the documents from the source collection
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}

		// Copy the documents to the destination collection
		destinationColl := destinationDB.Collection(collName)
		destinationColl.Drop(context.TODO())
		var documents []bson.M
		if err = cursor.All(ctx, &documents); err != nil {
			log.Fatal(err)
		}

		for _, doc := range documents {
			_, err := destinationColl.InsertOne(ctx, doc)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Database copy from", SOURCE_DB, "to", DESTINATION_DB, " complete!")
}
