package user

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "users"
	NEW_COLLECTION = "users"
)

type OldEntity struct {
	ID             string             `bson:"_id"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status"`
	Email          string             `bson:"email"`
	Name           string             `bson:"full_name"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	Role           string             `bson:"role"`
}

type NewEntity struct {
	ID        string             `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Status    string             `bson:"status"`
	Email     string             `bson:"email"`
	Name      string             `bson:"full_name"`
	OrgID     primitive.ObjectID `bson:"org_id"`
	Role      string             `bson:"role"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	newEntity := NewEntity{
		ID:        old.ID,
		CreatedAt: old.CreatedAt,
		UpdatedAt: old.UpdatedAt,
		Status:    old.Status,
		Email:     old.Email,
		Name:      old.Name,
		OrgID:     old.OrganizationID,
		Role:      old.Role,
	}
	return newEntity
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating User...")
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
		log.Println("Migrating user ", result.ID)
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
