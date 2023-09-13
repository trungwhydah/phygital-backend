package repository

import (
	"context"

	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TagRepository struct
type TagRepository struct {
	dbMongo *mongo.Database
}

// NewTagRepository create repository
func NewTagRepository(dbMongo *mongo.Database) *TagRepository {
	return &TagRepository{dbMongo: dbMongo}
}

func (r *TagRepository) CreateTag(tag *entity.Tag) (*entity.Tag, error) {
	result, err := r.dbMongo.Collection(tag.CollectionName()).InsertOne(context.TODO(), &tag)
	if err != nil {
		return nil, err
	}
	tag.ID = result.InsertedID.(primitive.ObjectID)

	return tag, nil
}

func (r *TagRepository) CheckExistedTag(tag *entity.Tag) (bool, error) {
	filter := bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{Key: "tag_id", Value: tag.TagID},
				},
				bson.D{
					{Key: "hardware_id", Value: tag.HardwareID},
				},
			},
		},
	}
	count, err := r.dbMongo.Collection(tag.CollectionName()).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (r *TagRepository) GetTagNotMapped(tagMapped *[]string) (*[]entity.Tag, error) {
	filter := bson.D{{Key: "tag_id", Value: bson.D{{Key: "$nin", Value: *tagMapped}}}}
	cursor, err := r.dbMongo.Collection(entity.Tag{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var tags []entity.Tag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		return nil, err
	}

	return &tags, nil
}

func (r *TagRepository) GetTag(tagID string) (*entity.Tag, error) {
	filter := bson.D{{Key: "tag_id", Value: tagID}}
	var tag entity.Tag
	err := r.dbMongo.Collection(tag.CollectionName()).FindOne(context.TODO(), filter).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &tag, nil
}

// UpdateChipCounter
func (r *TagRepository) UpdateTagCounter(tagID *string, scanCounter *int) (bool, error) {
	isUpsert := true
	option := options.UpdateOptions{Upsert: &isUpsert}
	filter := bson.D{{Key: "tag_id", Value: *tagID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "scan_counter", Value: *scanCounter},
		}}}
	result, err := r.dbMongo.Collection(entity.Tag{}.CollectionName()).UpdateOne(
		context.TODO(),
		filter,
		update,
		&option)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount != 0, nil
}

// GetTagByHWID
func (r *TagRepository) GetTagByHWID(uid *string) (*entity.Tag, error) {
	filter := bson.D{{Key: "hardware_id", Value: *uid}}
	var tag entity.Tag
	err := r.dbMongo.Collection(tag.CollectionName()).FindOne(context.TODO(), filter).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &tag, nil
}
