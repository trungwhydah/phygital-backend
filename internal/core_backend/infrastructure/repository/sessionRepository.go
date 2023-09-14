package repository

import (
	"backend-service/internal/core_backend/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	dbMongo *mongo.Database
}

// NewSessionRepository create repository
func NewSessionRepository(dbMongo *mongo.Database) *SessionRepository {
	return &SessionRepository{dbMongo: dbMongo}
}

func (sr *SessionRepository) CreateSession(session *entity.Session) (*entity.Session, error) {
	result, err := sr.dbMongo.Collection(session.CollectionName()).InsertOne(context.TODO(), &session)
	if err != nil {
		return nil, err
	}
	session.SessionID = result.InsertedID.(primitive.ObjectID)

	return session, nil
}

// GetSessionWithID
func (sr *SessionRepository) GetSessionWithID(sessionID *string) (*entity.Session, error) {
	sID, err := primitive.ObjectIDFromHex(*sessionID)
	if err != nil {
		return nil, err
	}

	var session entity.Session
	err = sr.dbMongo.Collection(session.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: sID}}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}
