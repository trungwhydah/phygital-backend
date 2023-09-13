package webPage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "webpages"
	NEW_COLLECTION = "webpages"
)

type OldEntity struct {
	ID         primitive.ObjectID       `bson:"_id"`
	CreatedAt  time.Time                `bson:"created_at"`
	UpdatedAt  time.Time                `bson:"updated_at"`
	Name       string                   `bson:"name"`
	Type       string                   `bson:"template"`
	URLLink    string                   `bson:"url_link"`
	Attributes []map[string]interface{} `bson:"attributes"`
}

type NewEntity struct {
	ID         primitive.ObjectID       `bson:"_id"`
	CreatedAt  time.Time                `bson:"created_at"`
	UpdatedAt  time.Time                `bson:"updated_at"`
	Status     string                   `bson:"status"`
	Name       string                   `bson:"name"`
	Type       string                   `bson:"type"`
	URLLink    string                   `bson:"url_link"`
	Category   string                   `bson:"category"`
	Attributes []map[string]interface{} `bson:"attributes"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	newURLLink := old.Type
	if newURLLink == "home" {
		newURLLink = ""
	}
	newEntity := NewEntity{
		ID:         old.ID,
		CreatedAt:  old.CreatedAt,
		UpdatedAt:  old.UpdatedAt,
		Status:     "Active",
		Name:       old.Name,
		Type:       old.Type,
		URLLink:    newURLLink,
		Category:   "",
		Attributes: old.Attributes,
	}
	return newEntity
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Webpage...")
	sourceCollection := sourceDB.Collection(OLD_COLLECTION)
	destinationCollection := destinationDB.Collection(NEW_COLLECTION)
	destinationCollection.Drop(context.Background())

	// Retrieve data from the source collection
	cursor, err := sourceCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cursor.Close(context.Background())

	// Iterate over the retrieved documents
	for cursor.Next(context.Background()) {
		var result OldEntity
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migrating web page ", result.ID)
		newResult := ConvertToNewEntity(result)

		// Insert the document into the destination collection
		_, err = destinationCollection.InsertOne(context.Background(), newResult)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
