package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	SessionID primitive.ObjectID `bson:"session_id"`
	TagID     string             `bson:"tag_id"`
	StartAt   time.Time          `bson:"start_at"`
	ExpiredAt time.Time          `bson:"expire_at"`
}

// CollectionName Collection name of Session
func (Session) CollectionName() string {
	return "sessions"
}
