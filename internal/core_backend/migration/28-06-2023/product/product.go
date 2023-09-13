package product

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_PRODUCT_ITEMS_COLLECTION = "product_items"
	OLD_COLLECTION               = "products"
	NEW_COLLECTION               = "products"
)

type OldEntity struct {
	ID             primitive.ObjectID `bson:"_id"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status"`
	ProductName    string             `bson:"product_name"`
	FarmName       string             `bson:"farm_name"`
	Varietal       string             `bson:"varietal"`
	Process        string             `bson:"process"`
	URLLink        string             `bson:"url_link"`
	TotalItem      int                `bson:"total_item"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	RatingScore    float64            `bson:"rating_score"`
}

type NewEntity struct {
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
	RatingScore float64            `bson:"rating_score"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	newEntity := NewEntity{
		ID:          old.ID,
		CreatedAt:   old.CreatedAt,
		UpdatedAt:   old.UpdatedAt,
		Status:      "Active",
		ProductName: old.ProductName,
		FarmName:    old.FarmName,
		Varietal:    old.Varietal,
		Process:     old.Process,
		URLLink:     old.URLLink,
		TotalItem:   old.TotalItem,
		OrgID:       old.OrganizationID,
		RatingScore: old.RatingScore,
	}
	return newEntity
}

func GetProductOrganizationMapping(productItemCollection *mongo.Collection) (map[string]string, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"product_id": bson.M{"$ne": primitive.NilObjectID},
			},
		},
		{
			"$group": bson.M{
				"_id":             "$product_id",
				"organization_id": bson.M{"$first": "$organization_id"},
			},
		},
	}

	cursor, err := productItemCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	mapping := make(map[string]string)
	for cursor.Next(context.TODO()) {
		var result struct {
			ProductID      primitive.ObjectID `bson:"_id"`
			OrganizationID string             `bson:"organization_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		mapping[result.ProductID.Hex()] = result.OrganizationID
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return mapping, nil
}

func GetOrganizationID(productID string, mapping map[string]string) primitive.ObjectID {
	var organizationID string
	if productID == "64a43a585ad8ec093255e8c3" {
		organizationID = "64a43a2f5ad8ec093255e8c2"
	} else {
		var exists bool
		organizationID, exists = mapping[productID]
		if !exists {
			log.Fatal(errors.New("Mapping doesn't exists"))
		}
	}

	orgIDHex, _ := primitive.ObjectIDFromHex(organizationID)

	return orgIDHex

}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Product...")
	sourceProductItemCollection := sourceDB.Collection(OLD_PRODUCT_ITEMS_COLLECTION)
	orgProdMapping, err := GetProductOrganizationMapping(sourceProductItemCollection)
	if err != nil {
		log.Fatal(err.Error())
	}

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
		log.Println("Migrating product ", result.ID)
		newResult := ConvertToNewEntity(result)
		orgID := GetOrganizationID(newResult.ID.Hex(), orgProdMapping)

		newResult.OrgID = orgID

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
