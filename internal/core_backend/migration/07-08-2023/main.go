package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"backend-service/internal/core_backend/config"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status    string             `bson:"status,omitempty" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type WebPageBase struct {
	BaseModel `bson:"inline"`
	Name      string `bson:"name,omitempty"`
	URLLink   string `bson:"url_link"`
	Type      string `bson:"type,omitempty"`
}

type NewWebpage struct {
	WebPageBase `bson:"inline"`
	Temp        any                    `bson:"-"`
	Attributes  map[string]interface{} `bson:"attributes"`
	// Attributes  map[string]interface{} `bson:"attributes"`
}

func main() {
	log.Println("Starting migration...")
	config.LoadConfig()
	var DBName = config.C.Mongo.DatabaseName
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.C.Mongo.MongoURLDbString))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Disconnect(context.Background())

	mess, err := webpagesArrayToObject(client.Database(DBName))
	if err != nil {
		log.Fatalln("Mingrated Step 1 failed with error: ", err)
	}

	log.Println("Mingrated successfully with message: ", mess)
	log.Println("Finish!")
	log.Println("❤️️")
}

func webpagesArrayToObject(db *mongo.Database) (string, error) {
	cursor, err := db.Collection("webpages").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	var oldWebpages []entity.WebPage
	if err = cursor.All(context.TODO(), &oldWebpages); err != nil {
		return "Failed", err
	}

	for _, wp := range oldWebpages {
		newWepage := NewWebpage{
			WebPageBase: wp.WebPageBase,
		}
		var resMap = make(map[string]interface{})
		for _, att := range wp.Attributes {
			keys := reflect.ValueOf(att).MapKeys()
			for _, j := range keys {
				resMap[j.String()] = att[j.String()]
			}
		}
		newWepage.Attributes = resMap

		_, err = db.Collection("new_webpages").InsertOne(context.TODO(), newWepage)
		if err != nil {
			return "Failed at: " + newWepage.Name, err
		}
	}

	return "Successfully!", nil
}
