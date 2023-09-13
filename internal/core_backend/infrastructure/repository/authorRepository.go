package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/entity"
)

// AuthorRepository struct
type AuthorRepository struct {
	dbMongo *mongo.Database
}

// NewAuthorRepository create repository
func NewAuthorRepository(dbMongo *mongo.Database) *AuthorRepository {
	return &AuthorRepository{dbMongo: dbMongo}
}

// CreateAuthor - create a new author
func (r *AuthorRepository) CreateAuthor(author *entity.Author) (*entity.Author, error) {
	result, err := r.dbMongo.Collection(author.CollectionName()).InsertOne(context.TODO(), author)
	if err != nil {
		return nil, err
	}
	author.ID = result.InsertedID.(primitive.ObjectID)

	return author, nil
}

// GetListAuthor - get all authors
func (r *AuthorRepository) GetListAuthor() (*[]entity.Author, error) {
	cursor, err := r.dbMongo.Collection(entity.Author{}.CollectionName()).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var listAuthors []entity.Author
	if err := cursor.All(context.TODO(), &listAuthors); err != nil {
		return nil, err
	}

	return &listAuthors, nil
}

// GetAuthorByID - get author by ID
func (r *AuthorRepository) GetAuthorByID(authorID *string) (*entity.Author, error) {
	aID, err := primitive.ObjectIDFromHex(*authorID)
	if err != nil {
		return nil, err
	}

	var author entity.Author
	err = r.dbMongo.Collection(author.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: aID}}).Decode(&author)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &author, nil
}
