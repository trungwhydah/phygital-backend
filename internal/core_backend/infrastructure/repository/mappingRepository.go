package repository

import (
	"context"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MappingRepository struct
type MappingRepository struct {
	dbMongo *mongo.Database
}

// NewMappingRepository create repository
func NewMappingRepository(dbMongo *mongo.Database) *MappingRepository {
	return &MappingRepository{dbMongo: dbMongo}
}

// GetAllMapping
func (r *MappingRepository) GetAllMapping() (*[]entity.Mapping, error) {
	cursor, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var mappings []entity.Mapping
	if err = cursor.All(context.TODO(), &mappings); err != nil {
		return nil, err
	}

	return &mappings, nil
}

// GetAllMappingInOrg
func (r *MappingRepository) GetAllMappingInOrg(orgID *primitive.ObjectID) (*[]entity.Mapping, error) {
	cursor, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).Find(
		context.TODO(),
		bson.D{{Key: "org_id", Value: *orgID}},
	)
	if err != nil {
		return nil, err
	}

	var mappings []entity.Mapping
	if err = cursor.All(context.TODO(), &mappings); err != nil {
		return nil, err
	}

	return &mappings, nil
}

// UpsertMapping
func (r *MappingRepository) UpsertMapping(mapping *entity.Mapping) (bool, error) {
	isUpsert := true
	option := options.UpdateOptions{Upsert: &isUpsert}
	filter := bson.M{"tag_id": mapping.TagID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "product_item_id", Value: mapping.ProductItemID},
			{Key: "org_id", Value: mapping.OrganizationID},
			{Key: "external_url", Value: mapping.ExternalURL},
			{Key: "claimable", Value: mapping.Claimable},
		}}}

	result, err := r.dbMongo.Collection(mapping.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update,
		&option)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount != 0, nil
}

// GetMappingWithTagID
func (r *MappingRepository) GetMappingWithTagID(tagID *string) (*entity.Mapping, error) {
	var mapping entity.Mapping
	err := r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "tag_id", Value: *tagID}}).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &mapping, nil
}

func (r *MappingRepository) GetMappingWithProductItemID(productItemID *string) (*entity.Mapping, error) {
	pID, _ := primitive.ObjectIDFromHex(*productItemID)
	var mapping entity.Mapping
	err := r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "product_item_id", Value: pID}}).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &mapping, nil
}

func (r *MappingRepository) GetAllMappingForProduct(productID *string, orgID *string) (*[]entity.Mapping, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return nil, err
	}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", entity.ProductItem{}.CollectionName()}, {"localField", "product_item_id"}, {"foreignField", "_id"}, {"as", "product_item"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$product_item"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"product_item.product_id", pID}}}}
	cursor, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})
	if err != nil {
		return nil, err
	}

	var mappings []entity.Mapping
	if err = cursor.All(context.TODO(), &mappings); err != nil {
		return nil, err
	}

	oID, _ := primitive.ObjectIDFromHex(*orgID)
	filter := bson.M{
		"$and": bson.A{
			bson.M{"product_item_id": primitive.NilObjectID},
			bson.M{"org_id": oID},
		},
	}
	cursor, err = r.dbMongo.Collection(entity.Mapping{}.CollectionName()).Find(context.TODO(), filter)

	for cursor.Next(context.TODO()) {
		var mapping entity.Mapping
		err := cursor.Decode(&mapping)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}
	if err != nil {
		return nil, err
	}
	return &mappings, nil
}

func (r *MappingRepository) UpdateMapping(tagID *string, req *request.UpdateMappingRequest) (bool, error) {
	filter := bson.D{{Key: "tag_id", Value: *tagID}}
	update := bson.D{{Key: "$set", Value: *req}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedMapping entity.Mapping
	err := r.dbMongo.Collection(updatedMapping.CollectionName()).FindOneAndUpdate(context.TODO(), filter, update, option).Decode(&updatedMapping)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MappingRepository) Unmap(tagID *string) (bool, error) {
	filter := bson.D{{Key: "tag_id", Value: *tagID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "product_item_id", Value: primitive.NilObjectID},
		}}}

	_, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MappingRepository) GetMappingByDigitalAsset(digitalAssetID *string) (*entity.Mapping, error) {
	dID, err := primitive.ObjectIDFromHex(*digitalAssetID)
	if err != nil {
		return nil, err
	}

	var mapping entity.Mapping
	err = r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "digital_asset_id", Value: dID}}).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &mapping, nil
}
