package repository

import (
	"context"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepository struct
type ProductRepository struct {
	dbMongo *mongo.Database
}

// NewProductRepository create repository
func NewProductRepository(dbMongo *mongo.Database) *ProductRepository {
	return &ProductRepository{dbMongo: dbMongo}
}

func (r *ProductRepository) CreateProduct(product *entity.Product) (*entity.Product, error) {
	result, err := r.dbMongo.Collection(product.CollectionName()).InsertOne(context.TODO(), &product)
	if err != nil {
		return nil, err
	}
	product.ID = result.InsertedID.(primitive.ObjectID)

	return product, nil
}

func (r *ProductRepository) CheckExistedProduct(product *entity.Product) (bool, error) {
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "status", Value: common.StatusActive}},
				bson.D{{Key: "product_name", Value: product.ProductName}},
				bson.D{{Key: "org_id", Value: product.OrganizationID}},
			},
		},
	}
	count, err := r.dbMongo.Collection(product.CollectionName()).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (r *ProductRepository) GetAllProducts() (*[]entity.Product, error) {
	cursor, err := r.dbMongo.Collection(entity.Product{}.CollectionName()).Find(context.TODO(), bson.D{{Key: "status", Value: common.StatusActive}})
	if err != nil {
		return nil, err
	}

	var products []entity.Product
	err = cursor.All(context.TODO(), &products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}

func (r *ProductRepository) GetAllProductsInOrg(orgID *primitive.ObjectID) (*[]entity.Product, error) {
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "status", Value: common.StatusActive}},
				bson.D{{Key: "org_id", Value: *orgID}},
			},
		},
	}
	cursor, err := r.dbMongo.Collection(entity.Product{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var products []entity.Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}

	return &products, nil
}

func (r *ProductRepository) GetProductByID(productID *string) (*entity.Product, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return nil, err
	}

	var product entity.Product
	err = r.dbMongo.Collection(product.CollectionName()).FindOne(context.TODO(), bson.M{"_id": pID}).Decode(&product)
	if err != nil {
		return nil, err
	}
	product.ParseAttribute()

	return &product, nil
}

// UpdateProductDetail
func (r *ProductRepository) UpdateProductDetail(product *entity.Product, productID *string) (bool, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return false, err
	}

	// log.Fatalln(product)
	var newProduct entity.Product
	err = r.dbMongo.Collection(product.CollectionName()).FindOneAndReplace(context.TODO(), bson.D{{Key: "_id", Value: pID}}, product).Decode(&newProduct)
	if err != nil {
		return false, err
	}

	return newProduct.ID == product.ID, nil
}

func (r *ProductRepository) SoftDeleteProductByID(productID *string) (bool, error) {
	oID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return false, err
	}

	filter := bson.D{{Key: "_id", Value: oID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: common.StatusInactive},
		}}}

	result, err := r.dbMongo.Collection(entity.Product{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update)

	if err != nil {
		return false, err
	}

	return result.ModifiedCount != 0, nil
}

func (r *ProductRepository) GetOrganizationNameByProductId(productID *string) (string, error) {
	iID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return "", err
	}

	var product entity.Product
	err = r.dbMongo.Collection(product.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: iID}}).Decode(&product)
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

func (r *ProductRepository) UpdateProductTotalItems(productID *string, totalItems int) (bool, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		return false, err
	}

	filter := bson.D{{Key: "_id", Value: pID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "total_item", Value: totalItems},
		}}}

	_, err = r.dbMongo.Collection(entity.Product{}.CollectionName()).UpdateOne(
		context.TODO(),
		&filter,
		&update)

	if err != nil {
		return false, err
	}
	return true, nil
}

// GetProductForAuthor -
func (r *ProductRepository) GetProductForAuthor(authorID *string) (*[]entity.Product, error) {
	aID, err := primitive.ObjectIDFromHex(*authorID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "status", Value: common.StatusActive}},
				bson.D{{Key: "org_id", Value: aID}},
			},
		},
	}

	cursor, err := r.dbMongo.Collection(entity.Product{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var products []entity.Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}

	return &products, nil
}
