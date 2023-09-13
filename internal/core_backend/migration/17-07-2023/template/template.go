package template

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	OLD_COLLECTION = "templates"
	NEW_COLLECTION = "templates"
)

type OldEntity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Status          string             `bson:"status,omitempty"`
	Name            string             `bson:"name"`
	CreatedByUserID string             `bson:"created_by_user_id"`
	Category        string             `bson:"category"`
	Languages       []string           `bson:"languages"`
	Pages           []OldTemplatePages `bson:"pages"`
	Menu            []OldTemplateMenu  `bson:"menu"`
}

type OldTemplatePages struct {
	Type   string             `bson:"type"`
	URL    string             `bson:"url"`
	PageID primitive.ObjectID `bson:"page_id"`
}

type OldTemplateMenu struct {
	Title TemplateMenuTitle `bson:"title"`
	URL   string            `bson:"url"`
}

type TemplateMenuTitle struct {
	VI string `bson:"vi"`
	EN string `bson:"en"`
}

type NewEntity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Status          string             `bson:"status,omitempty"`
	Name            string             `bson:"name"`
	CreatedByUserID string             `bson:"created_by_user_id"`
	Category        string             `bson:"category"`
	Languages       []string           `bson:"languages"`
	Pages           []NewTemplatePages `bson:"pages"`
	Menu            []NewTemplateMenu  `bson:"menu"`
}

type NewTemplatePages struct {
	PageID primitive.ObjectID `bson:"page_id"`
}

type NewTemplateMenu struct {
	Title  TemplateMenuTitle  `bson:"title"`
	PageID primitive.ObjectID `bson:"page_id"`
}

func ConvertToNewEntity(old OldEntity) NewEntity {
	var pages []NewTemplatePages
	for _, page := range old.Pages {
		pages = append(pages, NewTemplatePages{
			PageID: page.PageID,
		})
	}
	var menus []NewTemplateMenu
	for _, oldMenu := range old.Menu {
		var pageID primitive.ObjectID
		for _, oldPage := range old.Pages {
			if oldPage.URL == oldMenu.URL {
				pageID = oldPage.PageID
				break
			}
		}
		if len(pageID) == 0 {
			log.Println("Couldn't found corresponding template.menu pageID in template.pages")
		}
		menus = append(menus, NewTemplateMenu{
			Title:  oldMenu.Title,
			PageID: pageID,
		})
	}
	newEntity := NewEntity{
		ID:              old.ID,
		CreatedAt:       old.CreatedAt,
		UpdatedAt:       old.UpdatedAt,
		Status:          old.Status,
		Name:            old.Name,
		CreatedByUserID: old.CreatedByUserID,
		Category:        old.Category,
		Languages:       old.Languages,
		Pages:           pages,
		Menu:            menus,
	}

	return newEntity
}
func MigrateCollection(sourceDB *mongo.Database, destinationDB *mongo.Database) {
	log.Println("Start Migrating Template...")
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
		log.Println("Migrating template ", result.ID)
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
