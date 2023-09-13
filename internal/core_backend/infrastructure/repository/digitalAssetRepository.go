package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"
)

// DigitalAssetRepository struct
type DigitalAssetRepository struct {
	dbMongo *mongo.Database
}

// NewDigitalAssetRepository create repository
func NewDigitalAssetRepository(dbMongo *mongo.Database) *DigitalAssetRepository {
	return &DigitalAssetRepository{dbMongo: dbMongo}
}

func (r *DigitalAssetRepository) GetDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, error) {
	dID, err := primitive.ObjectIDFromHex(*collectionID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName()).Find(context.TODO(), bson.D{{Key: "collection_id", Value: dID}})
	if err != nil {
		return nil, err
	}
	var digitalAssets []entity.DigitalAsset
	if err = cursor.All(context.TODO(), &digitalAssets); err != nil {
		return nil, err
	}

	return &digitalAssets, nil
}

func (r *DigitalAssetRepository) GetActiveDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, error) {
	cID, err := primitive.ObjectIDFromHex(*collectionID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "status", Value: common.StatusActive}},
				bson.D{{Key: "collection_id", Value: cID}},
			},
		},
	}
	cursor, err := r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var digitalAssets []entity.DigitalAsset
	if err = cursor.All(context.TODO(), &digitalAssets); err != nil {
		return nil, err
	}

	return &digitalAssets, nil
}

// GetDigitalAssetByID
func (r *DigitalAssetRepository) GetDigitalAssetByID(daID *string) (*entity.DigitalAsset, error) {
	dID, err := primitive.ObjectIDFromHex(*daID)
	var da entity.DigitalAsset
	err = r.dbMongo.Collection(da.CollectionName()).FindOne(context.TODO(), bson.M{"_id": dID}).Decode(&da)
	if err != nil {
		return nil, err
	}

	return &da, nil
}

// GetDigitalAssetByTokenID
func (r *DigitalAssetRepository) GetDigitalAssetByTokenID(collectionID *string, tokenID *int) (*entity.DigitalAsset, error) {
	cID, err := primitive.ObjectIDFromHex(*collectionID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "status", Value: common.StatusActive}},
				bson.D{{Key: "collection_id", Value: cID}},
				bson.D{{Key: "token_id", Value: *tokenID}},
			},
		},
	}

	var asset entity.DigitalAsset
	err = r.dbMongo.Collection(asset.CollectionName()).FindOne(context.TODO(), filter).Decode(&asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// CreateDigitalAsset
func (r *DigitalAssetRepository) CreateDigitalAsset(da *entity.DigitalAsset) (*entity.DigitalAsset, error) {
	result, err := r.dbMongo.Collection(da.CollectionName()).InsertOne(context.TODO(), &da)
	if err != nil {
		return nil, err
	}
	da.ID = result.InsertedID.(primitive.ObjectID)

	return da, nil
}

// UpdateDigitalAsset
func (r *DigitalAssetRepository) UpdateDigitalAsset(da *entity.DigitalAsset) (bool, error) {
	_, err := r.dbMongo.Collection(da.CollectionName()).UpdateByID(
		context.TODO(),
		da.ID,
		bson.M{"$set": &da},
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *DigitalAssetRepository) UpdateDigitalAssetMetadata(daID *string, metadata *entity.Metadata) (bool, error) {
	daObjectID, err := primitive.ObjectIDFromHex(*daID)
	if err != nil {
		return false, err
	}
	filter := bson.D{{Key: "_id", Value: daObjectID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "metadata", Value: *metadata},
		}}}
	_, err = r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *DigitalAssetRepository) GetAllActiveDigitalAssets() (*[]entity.DigitalAsset, error) {
	cursor, err := r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName()).Find(context.TODO(), bson.D{{Key: "status", Value: common.StatusActive}})
	if err != nil {
		return nil, err
	}

	var digitalAssets []entity.DigitalAsset
	err = cursor.All(context.TODO(), &digitalAssets)
	if err != nil {
		return nil, err
	}

	return &digitalAssets, nil
}

func (r *DigitalAssetRepository) GetDigitalAssetsProductAggregate() (*[]entity.DigitalAssetProductAggregate, error) {
	coll := r.dbMongo.Collection(entity.DigitalAsset{}.CollectionName())
	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "mappings"},
					{"localField", "_id"},
					{"foreignField", "digital_asset_id"},
					{"as", "mapping"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$mapping"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"product_item_id", "$mapping.product_item_id"},
					{"org_id", "$mapping.org_id"},
				},
			},
		},
		bson.D{{"$unset", "mapping"}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "organizations"},
					{"localField", "org_id"},
					{"foreignField", "_id"},
					{"as", "organization"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$organization"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{{"$set", bson.D{{"org_tag_name", "$organization.org_tag_name"}}}},
		bson.D{
			{"$unset",
				bson.A{
					"organization",
					"org_id",
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "product_items"},
					{"localField", "product_item_id"},
					{"foreignField", "_id"},
					{"as", "product_item"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$product_item"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"product_id", "$product_item.product_id"},
					{"item_index", "$product_item.item_index"},
				},
			},
		},
		bson.D{{"$unset", "product_item"}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "products"},
					{"localField", "product_id"},
					{"foreignField", "_id"},
					{"as", "product"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$product"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{{"$unset", "product_id"}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var aggregations []entity.DigitalAssetProductAggregate
	if err = cursor.All(context.TODO(), &aggregations); err != nil {
		return nil, err
	}
	return &aggregations, nil
}
