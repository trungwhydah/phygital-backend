package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const WebPageCollectionName = "webpages"

// WebPageRepository struct
type WebPageRepository struct {
	dbMongo *mongo.Database
}

// WebPageRepository create repository
func NewWebPageRepository(dbMongo *mongo.Database) *WebPageRepository {
	return &WebPageRepository{dbMongo: dbMongo}
}

func (r *WebPageRepository) CreateWebPage(page *entity.WebPage) (*entity.WebPage, error) {
	result, err := r.dbMongo.Collection(page.CollectionName()).InsertOne(context.TODO(), &page)
	if err != nil {
		return nil, err
	}
	page.ID = result.InsertedID.(primitive.ObjectID)
	return page, nil
}

func (r *WebPageRepository) GetAllWebPages() (*[]entity.WebPage, error) {

	cursor, err := r.dbMongo.Collection(WebPageCollectionName).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var pages []entity.WebPage
	if err := cursor.All(context.TODO(), &pages); err != nil {
		return nil, err
	}

	return &pages, nil
}

func (r *WebPageRepository) GetWebPage(Id *string) (*entity.WebPage, error) {
	pageId, _ := primitive.ObjectIDFromHex(*Id)
	var webpage entity.WebPage
	err := r.dbMongo.Collection(webpage.CollectionName()).FindOne(context.TODO(), bson.M{"_id": pageId}).Decode(&webpage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &webpage, nil
}

func (r *WebPageRepository) UpdateWebPage(page *entity.WebPage) (*entity.WebPage, error) {

	_, err := r.dbMongo.Collection(page.CollectionName()).UpdateByID(
		context.TODO(),
		page.ID,
		bson.M{"$set": &page},
	)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (r *WebPageRepository) DeleteWebPage(Id *string) (bool, error) {

	pageId, _ := primitive.ObjectIDFromHex(*Id)
	_, err := r.dbMongo.Collection(WebPageCollectionName).DeleteOne(context.TODO(), bson.M{"_id": pageId})

	if err != nil {
		return false, err
	}

	return true, nil
}
