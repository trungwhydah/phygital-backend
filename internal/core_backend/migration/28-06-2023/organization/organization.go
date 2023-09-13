package organization

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "organizations"
	NEW_COLLECTION = "organizations"
)

type OldEntity struct {
	ID               primitive.ObjectID `bson:"_id"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
	Status           string             `bson:"status"`
	OrganizationName string             `bson:"organization_name"`
	NameTag          string             `bson:"name_tag"`
	LogoURL          string             `bson:"logo_url"`
	OwnerID          primitive.ObjectID `bson:"owner_id"`
}

type NewEntity struct {
	ID               primitive.ObjectID `bson:"_id"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
	Status           string             `bson:"status"`
	OrganizationName string             `bson:"org_name"`
	NameTag          string             `bson:"org_tag_name"`
	LogoURL          string             `bson:"org_logo_url"`
	OwnerID          string             `bson:"owner_id"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	newEntity := NewEntity{
		ID:               old.ID,
		CreatedAt:        old.CreatedAt,
		UpdatedAt:        old.UpdatedAt,
		Status:           "Active",
		OrganizationName: old.OrganizationName,
		NameTag:          old.NameTag,
		LogoURL:          old.LogoURL,
	}

	if old.OwnerID.IsZero() {
		newEntity.OwnerID = ""
	} else {
		newEntity.OwnerID = old.OwnerID.Hex()
	}

	return newEntity
}
func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Organization...")
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
		log.Println("Migrating organization ", result.ID)
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
