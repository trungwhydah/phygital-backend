package auto_sync_total_item

import (
	"context"
	"log"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID        primitive.ObjectID `bson:"_id"`
	TotalItem int                `bson:"total_item"`
}

func SyncTotalItems(database *mongo.Database) {
	log.Println("For each product, if #product-items < total_item, add new ones to meet total_item.\nIf not, reset total_item to #product-items")
	pItemCol := database.Collection("product_items")
	productCol := database.Collection("products")
	cursor, err := productCol.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product Product
		err := cursor.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		count, err := pItemCol.CountDocuments(context.TODO(), bson.D{{Key: "product_id", Value: product.ID}})
		if err != nil {
			log.Fatal(err)
		}
		if product.TotalItem < int(count) {
			log.Println("Product", product.ID.Hex(), " have total_item < real number of items")
			filter := bson.D{{Key: "_id", Value: product.ID}}
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "total_item", Value: int(count)},
				}}}
			_, err = productCol.UpdateOne(
				context.TODO(),
				&filter,
				&update)
		} else if product.TotalItem > int(count) {
			log.Println("Product", product.ID.Hex(), " have total_item > real number of items")
			for i := 0; i < product.TotalItem-int(count); i++ {
				log.Println("Adding product item with item_index", int(count)+i+1)
				item := &entity.ProductItem{
					ProductID: product.ID,
					BaseModel: entity.BaseModel{
						Status: common.StatusActive,
					},
					ItemIndex: int(count) + i + 1,
				}
				item.SetTime()
				pItemCol.InsertOne(context.TODO(), item)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
