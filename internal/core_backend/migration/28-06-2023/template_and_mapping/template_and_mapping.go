package template_and_mapping

import (
	"context"
	"log"
	"time"

	productMigration "backend-service/internal/core_backend/migration/28-06-2023/product"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Status       string             `bson:"status,omitempty"`
	ProductName  string             `bson:"product_name"`
	FarmName     string             `bson:"farm_name"`
	Varietal     string             `bson:"varietal"`
	Process      string             `bson:"process"`
	URLLink      string             `bson:"url_link"`
	TotalItem    int                `bson:"total_item"`
	RatingScore  float64            `bson:"rating_score"`
	ProductStory ProductStory       `bson:"product_story,omitempty"`
}

type ProductStory struct {
	Homepage    primitive.ObjectID            `bson:"homepage,omitempty"`
	SubPage     map[string]primitive.ObjectID `bson:"subpage,omitempty"`
	IsAstronaut bool                          `bson:"isAstronaut,omitempty"`
}

type Template struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Status          string             `bson:"status,omitempty"`
	Name            string             `bson:"name"`
	CreatedByUserID string             `bson:"created_by_user_id"`
	Category        string             `bson:"category"`
	Languages       []string           `bson:"languages"`
	Pages           []TemplatePages    `bson:"pages"`
	Menu            []TemplateMenu     `bson:"menu"`
}

type TemplatePages struct {
	Type   string             `bson:"type"`
	URL    string             `bson:"url"`
	PageID primitive.ObjectID `bson:"page_id"`
}

type TemplateMenu struct {
	Title TemplateMenuTitle `bson:"title"`
	URL   string            `bson:"url"`
}

type TemplateMenuTitle struct {
	VI string `bson:"vi"`
	EN string `bson:"en"`
}
type Mapping struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status,omitempty"`
	ProductItemID  primitive.ObjectID `bson:"product_item_id"`
	TagID          string             `bson:"tag_id"`
	TemplateID     primitive.ObjectID `bson:"template_id"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	ExternalURL    string             `bson:"external_url"`
	Claimable      bool               `bson:"claimable"`
	OwnerID        string             `bson:"owner_id"`
}

type ProductItem struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	Status         string             `bson:"status,omitempty"`
	ProductID      string             `bson:"product_id"`
	ChipID         string             `bson:"chip_id"`
	Claimable      bool               `bson:"claimable"`
	ExternalURL    string             `bson:"external_url"`
	OwnerID        string             `bson:"owner_id"`
	OrganizationID string             `bson:"organization_id"`
}

func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Template And Mapping...")
	// Access the "products" and "templates" collections
	productsCollection := sourceDB.Collection("products")
	productItemsCollection := sourceDB.Collection("product_items")
	templatesCollection := destinationDB.Collection("templates")
	mappingsCollection := destinationDB.Collection("mappings")
	templatesCollection.Drop(context.Background())
	mappingsCollection.Drop(context.Background())

	productWithOrg, err := productMigration.GetProductOrganizationMapping(productItemsCollection)
	if err != nil {
		log.Fatal(err)
	}
	// Retrieve the product from the "products" collection
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
		log.Println("Using product", product.ID.Hex(), "to migrate templates")
		// Retrieve the organization name using the "organization_id" from the "organizations" collection
		// Replace "organizations" with the actual name of your organizations collection
		var organization struct {
			Name string `bson:"organization_name"`
		}
		orgID := productMigration.GetOrganizationID(product.ID.Hex(), productWithOrg)
		err = sourceDB.Collection("organizations").FindOne(context.Background(), bson.M{"_id": orgID}).Decode(&organization)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new template document
		template := Template{
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Status:          "Active",
			CreatedByUserID: "",
			Name:            organization.Name + " - " + product.ProductName,
		}
		// Add the homepage to the "pages" array
		template.Pages = append(template.Pages, TemplatePages{
			Type:   "home",
			URL:    "",
			PageID: product.ProductStory.Homepage,
		})

		// Add the homepage to the "menu" array
		template.Menu = append(template.Menu, TemplateMenu{
			Title: TemplateMenuTitle{
				VI: "Trang chủ",
				EN: "Home",
			},
			URL: "",
		})

		// Map the "subpage" objects to the "pages" array
		for subpageKey, subpageID := range product.ProductStory.SubPage {
			template.Pages = append(template.Pages, TemplatePages{
				Type:   subpageKey,
				URL:    subpageKey,
				PageID: subpageID,
			})
			var titleVI string
			switch subpageKey {
			case "story":
				titleVI = "Câu chuyện"
			case "about":
				titleVI = "Giới thiệu"
			case "tutorial":
				titleVI = "Hướng dẫn"
			case "letter":
				titleVI = "Thư ngỏ"
			default:
				titleVI = subpageKey
			}
			template.Menu = append(template.Menu, TemplateMenu{
				Title: TemplateMenuTitle{
					VI: titleVI,
					EN: subpageKey,
				},
				URL: subpageKey,
			})
		}

		// Save the template document in the "templates" collection
		insertResult, err := templatesCollection.InsertOne(context.Background(), template)
		if err != nil {
			log.Fatal(err)
		}

		// Print the ID of the inserted template document
		log.Println("Inserted template ID:", insertResult.InsertedID)

		// Iterate through the product items for the current product
		productItemsCursor, err := productItemsCollection.Find(context.Background(), bson.M{"product_id": product.ID})
		if err != nil {
			log.Fatal(err)
		}
		defer productItemsCursor.Close(context.Background())

		// Iterate through the product items
		for productItemsCursor.Next(context.Background()) {
			var productItem ProductItem
			err := productItemsCursor.Decode(&productItem)
			if err != nil {
				log.Fatal(err)
			}

			// Create a new mapping document
			mapping := Mapping{
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Status:         "Active",
				ProductItemID:  productItem.ID,
				TagID:          productItem.ChipID,
				TemplateID:     insertResult.InsertedID.(primitive.ObjectID),
				OrganizationID: orgID,
				ExternalURL:    productItem.ExternalURL, // Set the external URL if applicable
				Claimable:      productItem.Claimable,
				OwnerID:        productItem.OwnerID,
			}

			// Save the mapping document in the "mappings" collection
			mp, err := mappingsCollection.InsertOne(context.Background(), mapping)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Inserted mapping ID:", mp.InsertedID)
		}
	}
	// Iterate through the product items for the current product
	productItemsCursor, err := productItemsCollection.Find(context.Background(), bson.M{"product_id": primitive.NilObjectID})
	if err != nil {
		log.Fatal(err)
	}
	defer productItemsCursor.Close(context.Background())

	// Iterate through the product items
	for productItemsCursor.Next(context.Background()) {
		var productItem ProductItem
		err := productItemsCursor.Decode(&productItem)
		if err != nil {
			log.Fatal(err)
		}
		orgID, _ := primitive.ObjectIDFromHex(productItem.OrganizationID)

		// Create a new mapping document
		mapping := Mapping{
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Status:         "Active",
			ProductItemID:  primitive.NilObjectID,
			TagID:          productItem.ChipID,
			TemplateID:     primitive.NilObjectID,
			OrganizationID: orgID,
			ExternalURL:    productItem.ExternalURL,
			Claimable:      productItem.Claimable,
			OwnerID:        productItem.OwnerID,
		}

		// Save the mapping document in the "mappings" collection
		mp, err := mappingsCollection.InsertOne(context.Background(), mapping)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted mapping ID:", mp.InsertedID)
	}
}
