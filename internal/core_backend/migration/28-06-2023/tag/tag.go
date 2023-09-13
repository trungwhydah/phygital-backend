package tag

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "chips"
	NEW_COLLECTION = "tags"
)

type OldEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	ChipID      string             `bson:"chip_id"`
	BatchID     string             `bson:"bactch_id"`
	ChipNumber  string             `bson:"chip_number"`
	ScanCounter int                `bson:"scan_counter"`
	ProductID   string             `bson:"product_id"`
	HardwareID  string             `bson:"hardware_id"`
}

type NewEntity struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status,omitempty"`
	HardwareID     string             `bson:"hardware_id"`
	TagID          string             `bson:"tag_id"`
	TagType        string             `bson:"tag_type"` // chip or qr
	EncryptMode    string             `bson:"encrypt_mode"`
	RawData        string             `bson:"raw_data"`
	ScanCounter    int                `bson:"scan_counter"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
}

func GetOrganizationIDFromChipID(sourceDB *mongo.Database, chipID string) primitive.ObjectID {
	productItemsColl := "product_items"

	var resultString struct {
		OrganizationID string `bson:"organization_id"`
	}
	err := sourceDB.Collection(productItemsColl).FindOne(context.Background(), bson.M{"chip_id": chipID}).Decode(&resultString)
	if err != nil {
		log.Fatal(err)
	}

	if resultString.OrganizationID == "" {
		return primitive.NilObjectID
	} else {
		orgID, err := primitive.ObjectIDFromHex(resultString.OrganizationID)
		if err != nil {
			log.Fatal(err)
		}
		return orgID
	}
}
func ConvertOldToNew(sourceDB *mongo.Database, old OldEntity) NewEntity {
	org_id := GetOrganizationIDFromChipID(sourceDB, old.ChipID)
	return NewEntity{
		ID:             old.ID,
		CreatedAt:      old.CreatedAt,
		UpdatedAt:      old.UpdatedAt,
		Status:         "Active",
		HardwareID:     old.HardwareID,
		TagID:          old.ChipID,
		TagType:        "",
		EncryptMode:    "",
		RawData:        "",
		ScanCounter:    old.ScanCounter,
		OrganizationID: org_id,
	}
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Tag...")
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
		log.Println("Migrating tag", result.ChipID)
		newResult := ConvertOldToNew(sourceDB, result)
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
