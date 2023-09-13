package move_template_to_product

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewProduct struct {
	ID          primitive.ObjectID `bson:"_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Status      string             `bson:"status"`
	ProductName string             `bson:"product_name"`
	FarmName    string             `bson:"farm_name"`
	Varietal    string             `bson:"varietal"`
	Process     string             `bson:"process"`
	URLLink     string             `bson:"url_link"`
	TotalItem   int                `bson:"total_item"`
	OrgID       primitive.ObjectID `bson:"org_id"`
	TemplateID  primitive.ObjectID `bson:"template_id"`
	RatingScore float64            `bson:"rating_score"`
}

type ProductItem struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Status    string             `bson:"status"`
	ProductID primitive.ObjectID `bson:"product_id"`
	OwnerID   string             `bson:"owner_id"`
	TotalLike int                `bson:"total_like"`
	ItemIndex int                `bson:"item_index"`
}

type OldMapping struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status,omitempty"`
	ProductItemID  primitive.ObjectID `bson:"product_item_id"`
	TagID          string             `bson:"tag_id"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	TemplateID     primitive.ObjectID `bson:"template_id"`
	ExternalURL    string             `bson:"external_url"`
	Claimable      bool               `bson:"claimable"`
	OwnerID        string             `bson:"owner_id"`
}
type NewMapping struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status,omitempty"`
	ProductItemID  primitive.ObjectID `bson:"product_item_id"`
	TagID          string             `bson:"tag_id"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	ExternalURL    string             `bson:"external_url"`
	Claimable      bool               `bson:"claimable"`
	OwnerID        string             `bson:"owner_id"`
}

func GetTemplateIDOfProduct(productID primitive.ObjectID, sourceDB *mongo.Database, destinationDB *mongo.Database) primitive.ObjectID {
	productItemColl := sourceDB.Collection("product_items")
	mappingsColl := sourceDB.Collection("mappings")
	cursor, err := productItemColl.Find(context.Background(), bson.M{"product_id": productID})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cursor.Close(context.Background())

	resultTemplateID := primitive.NilObjectID
	// Iterate over the retrieved documents
	for cursor.Next(context.Background()) {
		var productItem ProductItem
		err := cursor.Decode(&productItem)
		if err != nil {
			log.Fatal(err)
		}
		// Define the two filters
		filter1 := bson.M{"template_id": bson.M{"$ne": nil}}
		filter2 := bson.M{"product_item_id": productItem.ID}

		// Combine the filters using the $and operator
		var mapping OldMapping
		filter := bson.M{"$and": []bson.M{filter1, filter2}}
		err = mappingsColl.FindOne(context.Background(), filter).Decode(&mapping)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			log.Fatal("Error get template id of product (84): " + err.Error())
			continue
		}
		if resultTemplateID == primitive.NilObjectID {
			resultTemplateID = mapping.TemplateID
		} else {
			if resultTemplateID != mapping.TemplateID {
				log.Fatal("Product have > 1 template_id in their product_items (template_ids: " + resultTemplateID.Hex() + ", " + mapping.TemplateID.Hex() + ")")
				return primitive.NilObjectID
			}
		}
	}
	return resultTemplateID
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Template ID From Mapping To Product...")
	sourceProductCollection := sourceDB.Collection("products")
	destinationProductCollection := destinationDB.Collection("products")
	sourceMappingColl := sourceDB.Collection("mappings")
	destMappingColl := destinationDB.Collection("mappings")
	destinationProductCollection.Drop(context.Background())
	destMappingColl.Drop(context.Background())

	// Retrieve data from the source collection
	cursor, err := sourceProductCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cursor.Close(context.Background())

	// Iterate over the retrieved documents
	for cursor.Next(context.Background()) {
		var result NewProduct
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migrating product ", result.ID)
		tID := GetTemplateIDOfProduct(result.ID, sourceDB, destinationDB)

		result.TemplateID = tID

		// Insert the document into the destination collection
		_, err = destinationProductCollection.InsertOne(context.Background(), result)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Retrieve data from the source collection
	cursor, err = sourceMappingColl.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cursor.Close(context.Background())

	// Iterate over the retrieved documents
	for cursor.Next(context.Background()) {
		var result NewMapping
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migrating mapping ", result.ID)

		// Insert the document into the destination collection
		_, err = destMappingColl.InsertOne(context.Background(), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
