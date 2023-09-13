package verification

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "verifications"
	NEW_COLLECTION = "verifications"
)

type OldEntity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	TagID     string             `bson:"chip"`
	IsValid   bool               `bson:"is_valid"`
	Nonce     int                `bson:"nonce"`
}

type NewEntity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Status    string             `bson:"status"`
	TagID     string             `bson:"tag_id"`
	IsValid   bool               `bson:"is_valid"`
	Nonce     int                `bson:"nonce"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	newEntity := NewEntity{
		ID:        old.ID,
		CreatedAt: old.CreatedAt,
		UpdatedAt: old.UpdatedAt,
		Status:    "Active",
		TagID:     old.TagID,
		IsValid:   old.IsValid,
		Nonce:     old.Nonce,
	}
	return newEntity
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Verification...")
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
		log.Println("Migrating verification ", result.ID)
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
