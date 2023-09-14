package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Status    string             `bson:"status,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type OldProduct struct {
	BaseModel      `bson:"inline"`
	ProductName    string             `bson:"product_name"`
	Origin         string             `bson:"origin"`
	FarmName       string             `bson:"farm_name"`
	Varietal       string             `bson:"varietal"`
	Process        string             `bson:"process"`
	URLLink        string             `bson:"url_link"`
	TotalItem      int                `bson:"total_item"`
	TemplateID     primitive.ObjectID `bson:"template_id"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	RatingScore    float64            `bson:"rating_score"`
}

func main() {
	log.Println("Starting migration...")
	config.LoadConfig()
	var OldDBName = config.C.Mongo.DatabaseName
	var NewDBName = OldDBName + "_Migrate_2"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.C.Mongo.MongoURLDbString))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Disconnect(context.Background())

	// Initialize new database
	mess, err := copyCollectionToNewDatabase(OldDBName, NewDBName, []string{"webpages"})
	if err != nil {
		log.Fatalln("Error while init new Database: ", err.Error())
	}
	log.Println("Initialized new Database: ", mess)
	mess, isSuccess := migrateWebpageToProduct(client.Database(OldDBName), client.Database(NewDBName))
	if !isSuccess {
		log.Fatalln("Error while init new Database: ", err.Error())
	}
	log.Println("Mingrated successfully with message: ", mess)
	log.Println("Finish!")
	log.Println("❤️️")
}

func migrateWebpageToProduct(oldDatabase, newDatabase *mongo.Database) (string, bool) {
	cursor, err := oldDatabase.Collection("products").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer cursor.Close(context.Background())
	var oldProductList []OldProduct
	if err = cursor.All(context.TODO(), &oldProductList); err != nil {
		return err.Error(), false
	}

	for index, oldProduct := range oldProductList {
		log.Println("Start migrating product number:  " + strconv.Itoa(index+1))
		var template entity.Template
		err := oldDatabase.Collection("templates").FindOne(context.TODO(), bson.M{"_id": oldProduct.TemplateID}).Decode(&template)
		if err != nil {
			log.Println("Error: " + err.Error() + " while getting template!")
			return err.Error(), false
		}

		newProduct := entity.Product{
			BaseModel: entity.BaseModel{
				ID:        oldProduct.ID,
				Status:    oldProduct.Status,
				CreatedAt: oldProduct.CreatedAt,
				UpdatedAt: oldProduct.UpdatedAt,
			},
			ProductName:    oldProduct.ProductName,
			URLLink:        oldProduct.URLLink,
			TotalItem:      oldProduct.TotalItem,
			TemplateID:     oldProduct.TemplateID,
			OrganizationID: oldProduct.OrganizationID,
			RatingScore:    oldProduct.RatingScore,
		}

		var organization entity.Organization
		err = oldDatabase.Collection("organizations").FindOne(context.TODO(), bson.M{"_id": oldProduct.OrganizationID}).Decode(&organization)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("Error: " + err.Error() + " while getting organization for product number " + strconv.Itoa(index))
			return err.Error(), false
		}

		mess := fmt.Sprintf("Migrate product ID: %s for Organization: %s!", oldProduct.ID.Hex(), organization.NameTag)
		log.Println(mess)
		switch organization.NameTag {
		case "lej":
			var (
				// pageCounter int
				translation = make(map[string]interface{})
				viMap       = make(map[string]interface{})
				enMap       = make(map[string]interface{})
			)

			newProduct.Type = "coffee"
			coffeeAttribute := entity.AttributeCoffee{
				FarmName: oldProduct.FarmName,
				Varietal: oldProduct.Varietal,
				Process:  oldProduct.Process,
			}

			for _, page := range template.Pages {
				var pageAttribute entity.WebPage
				err = oldDatabase.Collection("webpages").FindOne(context.TODO(), bson.M{"_id": page.PageID}).Decode(&pageAttribute)
				if err != nil {
					log.Println("Error: " + err.Error() + " while getting page data for product number " + strconv.Itoa(index))
					return err.Error(), false
				}
				var mapAttributeVI = make(map[string]interface{})
				var mapAttributeEN = make(map[string]interface{})
				log.Println("product: ", oldProduct.ID.Hex(), "page type: ", pageAttribute.Type, " and page id: ", pageAttribute.ID.Hex())
				switch pageAttribute.Type {
				case "home":
					for _, att := range pageAttribute.Attributes {
						if att["vi"] != nil {
							mapAttributeVI = att["vi"].(map[string]interface{})
						}
						if att["en"] != nil {
							mapAttributeEN = att["en"].(map[string]interface{})
						}
					}
					imageAttribute := fmt.Sprintf("%v", mapAttributeVI["image"])
					newProduct.Image = entity.Media{
						Type: "image",
						URL:  imageAttribute,
					}
					viMap["description"] = mapAttributeVI["description"]
					viMap["title"] = mapAttributeVI["title"]

					enMap["description"] = mapAttributeEN["description"]
					enMap["title"] = mapAttributeEN["title"]

					var mediaCount int
					for _, att := range pageAttribute.Attributes {
						if att["image"] != nil {
							newProduct.Image.URL = fmt.Sprintf("%v", att["image"])
							mediaCount++
						}

						if att["videoThumbnail"] != nil {
							newProduct.Video.ThumbnailURL = fmt.Sprintf("%v", att["videoThumbnail"])
							mediaCount++
						}

						if mediaCount == 2 {
							break
						}
					}
				case "story":
					coffeeAttribute.CountryVideo = entity.Media{Type: "video"}
					coffeeAttribute.FarmVideo = entity.Media{Type: "video"}
					coffeeAttribute.CountryImage = entity.Media{Type: "image"}
					for _, att := range pageAttribute.Attributes {
						if att["vi"] != nil {
							mapAttributeVI = att["vi"].(map[string]interface{})
						}
						if att["en"] != nil {
							mapAttributeEN = att["en"].(map[string]interface{})
						}

						if att["countryVideoUrl"] != nil {
							coffeeAttribute.CountryVideo.URL = fmt.Sprintf("%v", att["countryVideoUrl"])
						}

						if att["countryVideoThumbnailUrl"] != nil {
							coffeeAttribute.CountryVideo.ThumbnailURL = fmt.Sprintf("%v", att["countryVideoThumbnailUrl"])
						}

						if att["preProcessingAcidity"] != nil {
							coffeeAttribute.Acidity = fmt.Sprintf("%v", att["preProcessingAcidity"])
						}

						if att["preProcessingBitter"] != nil {
							coffeeAttribute.Bitter = fmt.Sprintf("%v", att["preProcessingBitter"])
						}

						if att["preProcessingSweet"] != nil {
							coffeeAttribute.Sweet = fmt.Sprintf("%v", att["preProcessingSweet"])
						}

						if att["countryUrl"] != nil {
							coffeeAttribute.CountryImage.URL = fmt.Sprintf("%v", att["countryUrl"])
						}

						if att["farmVideoUrl"] != nil {
							coffeeAttribute.FarmVideo.URL = fmt.Sprintf("%v", att["farmVideoUrl"])
						}
						if att["farmVideoThumbnailUrl"] != nil {
							coffeeAttribute.FarmVideo.ThumbnailURL = fmt.Sprintf("%v", att["farmVideoThumbnailUrl"])
						}
					}

					newProduct.Origin = fmt.Sprintf("%v", mapAttributeVI["countryName"])
					viMap["countryName"] = fmt.Sprintf("%v", mapAttributeVI["countryName"])
					viMap["fact1Title"] = fmt.Sprintf("%v", mapAttributeVI["fact1Title"])
					viMap["fact1Description"] = fmt.Sprintf("%v", mapAttributeVI["fact1Description"])
					viMap["fact2Title"] = fmt.Sprintf("%v", mapAttributeVI["fact2Title"])
					viMap["fact2Description"] = fmt.Sprintf("%v", mapAttributeVI["fact2Description"])
					viMap["fact3Description"] = fmt.Sprintf("%v", mapAttributeVI["fact3Description"])
					viMap["farmName"] = fmt.Sprintf("%v", mapAttributeVI["farmName"])
					viMap["farmShortDescription"] = fmt.Sprintf("%v", mapAttributeVI["farmShortDescription"])
					viMap["farmFullDescription"] = fmt.Sprintf("%v", mapAttributeVI["farmFullDescription"])
					viMap["preProcessingDescription"] = fmt.Sprintf("%v", mapAttributeVI["preProcessingDescription"])
					viMap["preProcessingName"] = fmt.Sprintf("%v", mapAttributeVI["preProcessingName"])

					enMap["countryName"] = fmt.Sprintf("%v", mapAttributeEN["countryName"])
					enMap["fact1Title"] = fmt.Sprintf("%v", mapAttributeEN["fact1Title"])
					enMap["fact1Description"] = fmt.Sprintf("%v", mapAttributeEN["fact1Description"])
					enMap["fact2Title"] = fmt.Sprintf("%v", mapAttributeEN["fact2Title"])
					enMap["fact2Description"] = fmt.Sprintf("%v", mapAttributeEN["fact2Description"])
					enMap["fact3Description"] = fmt.Sprintf("%v", mapAttributeEN["fact3Description"])
					enMap["farmName"] = fmt.Sprintf("%v", mapAttributeEN["farmName"])
					enMap["farmShortDescription"] = fmt.Sprintf("%v", mapAttributeEN["farmShortDescription"])
					enMap["farmFullDescription"] = fmt.Sprintf("%v", mapAttributeEN["farmFullDescription"])
					enMap["preProcessingDescription"] = fmt.Sprintf("%v", mapAttributeEN["preProcessingDescription"])
					enMap["preProcessingName"] = fmt.Sprintf("%v", mapAttributeEN["preProcessingName"])

					var farmAttCount int
					for _, att := range pageAttribute.Attributes {
						if att["farmHeight"] != nil {
							coffeeAttribute.FarmHeight = fmt.Sprintf("%v", att["farmHeight"])
							farmAttCount++
						}
						if att["farmArea"] != nil {
							coffeeAttribute.FarmArea = fmt.Sprintf("%v", att["farmArea"])
							farmAttCount++
						}
						if farmAttCount == 2 {
							break
						}
					}
				case "tutorial":
					var tutorialAttCount int
					for _, att := range pageAttribute.Attributes {
						if att["brewingTime"] != nil {
							coffeeAttribute.BrewingTime = fmt.Sprintf("%v", att["brewingTime"])
							tutorialAttCount++
						}
						if att["buyUrl"] != nil {
							coffeeAttribute.LinkBuyProduct = fmt.Sprintf("%v", att["buyUrl"])
							tutorialAttCount++
						}
						if tutorialAttCount == 2 {
							break
						}
					}
				}
			}
			translation["vi"] = viMap
			translation["en"] = enMap
			coffeeAttribute.Translation = translation
			newProduct.Attribute = coffeeAttribute
			addNewProduct(newDatabase, newProduct)
		case "da-non-nuoc":
			newProduct.Type = "sculpture"
			newProduct.ThreeDimension = entity.Media{Type: "3D"}
			var stoneAttribute entity.AttributeSculpture
			for _, page := range template.Pages {
				var pageAttribute entity.WebPage
				err = oldDatabase.Collection("webpages").FindOne(context.TODO(), bson.M{"_id": page.PageID}).Decode(&pageAttribute)
				if err != nil {
					log.Println("Error: " + err.Error() + " while getting page data for product number " + strconv.Itoa(index))
					return err.Error(), false
				}
				switch pageAttribute.Type {
				case "home":
					for _, att := range pageAttribute.Attributes {
						if att["stone3dModel"] != nil {
							newProduct.ThreeDimension.URL = fmt.Sprintf("%v", att["stone3dModel"])
						}
						if att["stone3dThumbnail"] != nil {
							newProduct.ThreeDimension.ThumbnailURL = fmt.Sprintf("%v", att["stone3dThumbnail"])
						}
						if att["stoneTime"] != nil {
							stoneAttribute.SculptureTime = fmt.Sprintf("%v", att["stoneTime"])
						}
					}
				case "craft_village":
					mapVillageName := pageAttribute.Attributes[3]["villageName"].(map[string]interface{})
					stoneAttribute.Village = entity.Village{
						Name: fmt.Sprintf("%v", mapVillageName["vi"]),
						LocationVideo: entity.Media{
							Type: "video",
							URL:  fmt.Sprintf("%v", pageAttribute.Attributes[7]["locationVideoUrl"]),
						},
					}
				case "craftsmen":
					var (
						craftsman         entity.Craftsman
						craftsmanAttCount int
					)

					for _, att := range pageAttribute.Attributes {
						if att["craftsmenName"] != nil {
							craftsman.Name = fmt.Sprintf("%v", att["craftsmenName"])
							craftsmanAttCount++
						}
						if att["experience"] != nil {
							craftsman.ExperienceYear = fmt.Sprintf("%v", att["experience"])
							craftsmanAttCount++
						}
						if att["artworkCount"] != nil {
							craftsman.ArtworksCount = fmt.Sprintf("%v", att["artworkCount"])
							craftsmanAttCount++
						}
						if att["phone"] != nil {
							craftsman.Phone = fmt.Sprintf("%v", att["phone"])
							craftsmanAttCount++
						}
						if att["email"] != nil {
							craftsman.Email = fmt.Sprintf("%v", att["email"])
							craftsmanAttCount++
						}
						if att["description"] != nil {
							craftsman.Description = fmt.Sprintf("%v", att["description"])
							craftsmanAttCount++
						}
						if att["avatar"] != nil {
							craftsman.Avatar = entity.Media{
								URL:  fmt.Sprintf("%v", att["avatar"]),
								Type: "image",
							}
							craftsmanAttCount++
						}

						if craftsmanAttCount == 7 {
							break
						}
					}
					stoneAttribute.Craftsman = craftsman
				case "stone_info":
					var (
						stoneInfo entity.Stone
						count     int
					)
					for _, att := range pageAttribute.Attributes {
						if att["stoneName"] != nil {
							stoneInfo.Name = fmt.Sprintf("%v", att["stoneName"])
							count++
						}
						if att["origin"] != nil {
							stoneInfo.Origin = fmt.Sprintf("%v", att["origin"])
							count++
						}
						if att["color"] != nil {
							stoneInfo.Color = fmt.Sprintf("%v", att["color"])
							count++
						}
						if att["rarity"] != nil {
							stoneInfo.Rarity = fmt.Sprintf("%v", att["rarity"])
							count++
						}
						if att["characteristic"] != nil {
							stoneInfo.Properties = fmt.Sprintf("%v", att["characteristic"])
							count++
						}
						if att["imageUrl"] != nil {
							stoneInfo.Image = entity.Media{
								URL:  fmt.Sprintf("%v", att["imageUrl"]),
								Type: "image",
							}
							count++
						}
						if count == 6 {
							break
						}
					}
					stoneAttribute.Stone = stoneInfo
				}
			}
			newProduct.Attribute = stoneAttribute
			addNewProduct(newDatabase, newProduct)
		case "astronaut":
			newProduct.Type = "astronaut"
			newProduct.Video = entity.Media{Type: "video"}
			newProduct.Image = entity.Media{Type: "image"}
			for _, page := range template.Pages {
				var pageAttribute entity.WebPage
				err = oldDatabase.Collection("webpages").FindOne(context.TODO(), bson.M{"_id": page.PageID}).Decode(&pageAttribute)
				if err != nil {
					log.Println("Error: " + err.Error() + " while getting page data for product number " + strconv.Itoa(index))
					return err.Error(), false
				}

				if pageAttribute.Type == "home" {
					for _, att := range pageAttribute.Attributes {
						if att["video"] != nil {
							newProduct.Video.URL = fmt.Sprintf("%v", att["video"])
						}

						if att["videoThumbnail"] != nil {
							newProduct.Video.ThumbnailURL = fmt.Sprintf("%v", att["videoThumbnail"])
						}

						if att["image"] != nil {
							newProduct.Image.URL = fmt.Sprintf("%v", att["image"])
						}
					}
				}
			}
			addNewProduct(newDatabase, newProduct)
		case "ortho":
			newProduct.Type = "ortho"
			addNewProduct(newDatabase, newProduct)
		default:
			log.Println("Not found organization for product number ", strconv.Itoa(index), " with ID: ", oldProduct.ID.Hex())
		}
	}

	return "Successfully migrated!", true
}

func copyCollectionToNewDatabase(OldDBName, NewDBName string, listCollection []string) (string, error) {
	DatabaseURLString := config.C.Mongo.MongoURLDbString
	sourceClient, err := mongo.NewClient(options.Client().ApplyURI(DatabaseURLString))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()
	err = sourceClient.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer sourceClient.Disconnect(ctx)

	// Connect to the destination MongoDB server
	destinationClient, err := mongo.NewClient(options.Client().ApplyURI(DatabaseURLString))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3600*time.Second)
	defer cancel()
	err = destinationClient.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer destinationClient.Disconnect(ctx)

	// Get the source and destination databases
	sourceDB := sourceClient.Database(OldDBName)
	destinationDB := destinationClient.Database(NewDBName)

	// Copy each collection from the source to the destination database
	for _, collName := range listCollection {
		log.Println("Migrating", collName, "...")
		coll := sourceDB.Collection(collName)

		// Read the documents from the source collection
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			log.Fatalln(err)
		}

		// Copy the documents to the destination collection
		destinationColl := destinationDB.Collection(collName)
		destinationColl.Drop(context.TODO())
		var documents []bson.M
		if err = cursor.All(ctx, &documents); err != nil {
			log.Fatalln(err)
		}

		for _, doc := range documents {
			_, err := destinationColl.InsertOne(ctx, doc)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return "Successfully initialized database from " + OldDBName + " to " + NewDBName, nil
}

func addNewProduct(db *mongo.Database, newProduct entity.Product) error {
	res, err := db.Collection("products").InsertOne(context.TODO(), newProduct)
	if err != nil {
		log.Println("Error: " + err.Error() + " while creating new product for product ID " + newProduct.ID.Hex())
		return err
	}
	if res.InsertedID.(primitive.ObjectID) == primitive.NilObjectID {
		log.Println("Inserted product ID " + newProduct.ID.Hex() + " with null Object ID")
	}

	log.Println("Successfully created new product for product ID " + newProduct.ID.Hex())
	log.Println("----------------------------------------------------")
	return nil
}
