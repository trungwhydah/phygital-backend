package repository

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-service/internal/core_backend/entity"
)

// ProductItemRepository struct
type ProductItemRepository struct {
	dbMongo *mongo.Database
}

// NewProductItemRepository create repository
func NewProductItemRepository(dbMongo *mongo.Database) *ProductItemRepository {
	return &ProductItemRepository{dbMongo: dbMongo}
}

func (r *ProductItemRepository) InsertProductItem(item *entity.ProductItem) (*entity.ProductItem, error) {
	result, err := r.dbMongo.Collection(item.CollectionName()).InsertOne(context.TODO(), item)
	if err != nil {
		return nil, err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (r *ProductItemRepository) CheckProductItemMapped(productItemID *string) (bool, error) {
	if productItemID == nil {
		return false, nil
	}
	productItemObjectID, err := primitive.ObjectIDFromHex(*productItemID)
	if productItemObjectID.IsZero() {
		return false, nil
	}
	filter := bson.D{{Key: "product_item_id", Value: productItemObjectID}}
	count, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (r *ProductItemRepository) GetAllProductItemByProductID(productID *primitive.ObjectID) (*[]entity.ProductItem, error) {
	cursor, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).Find(
		context.TODO(),
		bson.D{{Key: "product_id", Value: *productID}},
	)
	if err != nil {
		return nil, err
	}

	var productItems []entity.ProductItem
	if err = cursor.All(context.TODO(), &productItems); err != nil {
		return nil, err
	}

	return &productItems, nil

}

func (r *ProductItemRepository) GetAllProductItem() (*[]entity.ProductItem, error) {
	cursor, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var productItems []entity.ProductItem
	if err = cursor.All(context.TODO(), &productItems); err != nil {
		return nil, err
	}

	return &productItems, nil
}

func (r *ProductItemRepository) GetDetailProductItemByID(productItemID *string) (*entity.ProductItem, error) {
	iID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return nil, err
	}
	var item entity.ProductItem
	err = r.dbMongo.Collection(item.CollectionName()).FindOne(context.TODO(), bson.M{"_id": iID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func (r *ProductItemRepository) IsAbleToClaim(productItemID *string) (bool, error) {
	iID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return false, err
	}

	var mapping entity.Mapping
	err = r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), bson.M{"product_item_id": iID}).Decode(&mapping)
	if err != nil {
		return false, err
	}

	return mapping.Claimable, nil
}

func (r *ProductItemRepository) SetOwner(productItemID, ownerID *string) (bool, error) {
	iID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}

	isUpsert := true
	option := options.UpdateOptions{Upsert: &isUpsert}
	filter := bson.D{{Key: "product_item_id", Value: iID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "owner_id", Value: *ownerID},
			{Key: "claimable", Value: false},
		}}}
	result, err := r.dbMongo.Collection(entity.Mapping{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update,
		&option)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount != 0, nil
}

func (r *ProductItemRepository) GetDetailWithTagID(tagID *string) (*entity.ProductItem, error) {
	var mapping entity.Mapping
	err := r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), bson.M{"tag_id": tagID}).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	var item entity.ProductItem
	err = r.dbMongo.Collection(item.CollectionName()).FindOne(context.TODO(), bson.M{"_id": mapping.ProductItemID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// This function used to count limited edition - Will be refactored later
func (r *ProductItemRepository) CountLimitedEdition(productID, itemID string) (string, error) {
	pID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return "", err
	}

	filter := bson.D{{Key: "product_id", Value: pID}}
	cursor, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return "", err
	}

	var items []entity.ProductItem
	if err = cursor.All(context.TODO(), &items); err != nil {
		return "", err
	}

	var result string
	for index, item := range items {
		if item.ID.Hex() == itemID {
			result = strconv.Itoa(index + 1)
		}
	}

	return result, nil
}

func (r *ProductItemRepository) GetOrganizationNameByProductItemID(productItemID *string) (string, error) {
	iID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return "", err
	}

	var item entity.ProductItem

	err = r.dbMongo.Collection(item.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: iID}}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		return "", err
	}

	var product entity.Product
	err = r.dbMongo.Collection(product.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: item.ProductID}}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		return "", err
	}

	var org entity.Organization
	err = r.dbMongo.Collection(org.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: product.OrganizationID}}).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		return "", err
	}

	return org.OrganizationName, nil
}

func (r *ProductItemRepository) UpdateTotalLike(productItemID *string) (bool, error) {
	piID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return false, err
	}
	filter := bson.D{{Key: "_id", Value: piID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "total_like", Value: 1},
		}}}
	result, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount != 0, nil
}

func (r *ProductItemRepository) ToggleClaimable(productItemID *string) (bool, error) {
	iID, err := primitive.ObjectIDFromHex(*productItemID)
	if err != nil {
		return false, err
	}

	var mapping entity.Mapping

	filter := bson.D{{Key: "product_item_id", Value: iID}}
	err = r.dbMongo.Collection(mapping.CollectionName()).FindOne(context.TODO(), filter).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "claimable", Value: !mapping.Claimable},
		}}}

	resultUpdate, err := r.dbMongo.Collection(mapping.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update)
	if err != nil {
		return false, err
	}

	return resultUpdate.ModifiedCount+resultUpdate.UpsertedCount+resultUpdate.MatchedCount != 0, nil
}

func (r *ProductItemRepository) CountNumProductItems(productID *string) (int, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return 0, err
	}
	count, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).CountDocuments(context.TODO(), bson.D{{Key: "product_id", Value: pID}})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ProductItemRepository) GetProductItemsInOrg(orgID *string) (*[]entity.ProductItem, error) {
	oID, err := primitive.ObjectIDFromHex(*orgID)
	if err != nil {
		return nil, err
	}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", entity.Product{}.CollectionName()}, {"localField", "product_id"}, {"foreignField", "_id"}, {"as", "product_join"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$product_join"}, {"preserveNullAndEmptyArrays", false}}}}
	filterStage := bson.D{{"$match", bson.D{{"product_join.org_id", oID}}}}
	cursor, err := r.dbMongo.Collection(entity.ProductItem{}.CollectionName()).Aggregate(context.TODO(), mongo.Pipeline{lookupStage, unwindStage, filterStage})
	if err != nil {
		return nil, err
	}

	var pItems []entity.ProductItem
	if err = cursor.All(context.TODO(), &pItems); err != nil {
		return nil, err
	}

	return &pItems, nil
}

func (r *ProductItemRepository) GetProductItemProductOrgAggregate(pItemID *string) (*entity.ProductItemProductOrgAggregate, error) {
	pID, err := primitive.ObjectIDFromHex(*pItemID)
	if err != nil {
		return nil, err
	}
	coll := r.dbMongo.Collection(entity.ProductItem{}.CollectionName())
	cursor, err := coll.Aggregate(context.TODO(), bson.A{
		bson.D{{"$match", bson.D{{"_id", pID}}}},
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
		bson.D{
			{"$set",
				bson.D{
					{"org_id", "$product.org_id"},
				},
			},
		},
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
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var aggregations []entity.ProductItemProductOrgAggregate
	if err = cursor.All(context.TODO(), &aggregations); err != nil {
		return nil, err
	}
	if len(aggregations) != 1 {
		err = errors.New("There must be exactly 1 product item with its org found! But there are != 1!")
		return nil, err
	}
	return &aggregations[0], nil
}
