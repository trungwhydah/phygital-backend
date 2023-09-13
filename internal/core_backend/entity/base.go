package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"  binding:"omitempty"`
	Status    string             `bson:"status,omitempty" json:"status"  binding:"omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at" binding:"omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"  binding:"omitempty"`
}

func (bm *BaseModel) SetTime() *BaseModel {
	if bm.CreatedAt.IsZero() {
		bm.CreatedAt = time.Now()
	}
	bm.UpdatedAt = time.Now()

	return bm
}

func (bm *BaseModel) SetStatus(status string) *BaseModel {
	bm.Status = status
	return bm
}

func (bm *BaseModel) ClearModel() *BaseModel {
	bm.ID = primitive.NilObjectID

	return bm
}

type Media struct {
	URL          string `bson:"url" json:"url"`
	Type         string `bson:"type" json:"type"`
	ThumbnailURL string `bson:"thumbnail_url" json:"thumbnail_url"`
}
