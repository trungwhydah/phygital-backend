package productItem

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OLD_COLLECTION = "product_items"
	NEW_COLLECTION = "product_items"
)

type OldEntity struct {
	ID              primitive.ObjectID `bson:"_id"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	ProductID       primitive.ObjectID `bson:"product_id"`
	ProductVariable primitive.ObjectID `bson:"product_variable"`
	ChipID          string             `bson:"chip_id"`
	Claimable       bool               `bson:"claimable"`
	OwnerID         string             `bson:"owner_id"`
	OrganizationID  string             `bson:"organization_id"`
	ExternalURL     string             `bson:"external_url"`
	TotalLike       int                `bson:"total_like"`
}

type NewEntity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Status    string             `bson:"status"`
	ProductID primitive.ObjectID `bson:"product_id"`
	TotalLike int                `bson:"total_like"`
	ItemIndex int                `bson:"item_index"`
}

type Product struct {
	ID primitive.ObjectID `bson:"_id"`
}

func ConvertToNewEntity(old OldEntity, index int) NewEntity {
	newEntity := NewEntity{
		ID:        old.ID,
		CreatedAt: old.CreatedAt,
		UpdatedAt: old.UpdatedAt,
		Status:    "Active",
		ProductID: old.ProductID,
		TotalLike: old.TotalLike,
		ItemIndex: index,
	}
	return newEntity
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Product Item...")
	productsCollection := sourceDB.Collection("products")
	sourceCollection := sourceDB.Collection(OLD_COLLECTION)
	destinationCollection := destinationDB.Collection(NEW_COLLECTION)
	destinationCollection.Drop(context.Background())

	cursor, err := productsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	// Iterate through the products
	for cursor.Next(context.Background()) {
		var product Product
		err := cursor.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Using product", product.ID.Hex(), "to migrate product items")
		sortOptions := options.Find().SetSort(bson.D{{"_id", 1}})
		// Iterate through the product items for the current product
		productItemsCursor, err := sourceCollection.Find(context.Background(), bson.M{"product_id": product.ID}, sortOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer productItemsCursor.Close(context.Background())

		index := 1
		// Iterate through the product items
		for productItemsCursor.Next(context.Background()) {
			var result OldEntity
			err := productItemsCursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			if result.ProductID == primitive.NilObjectID {
				log.Println("Product item doesn't have product id. Skipped...", result.ID)
				continue
			}
			log.Println("Migrating product item ", result.ID)
			newResult := ConvertToNewEntity(result, index)
			index += 1

			// Insert the document into the destination collection
			_, err = destinationCollection.InsertOne(context.Background(), newResult)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
