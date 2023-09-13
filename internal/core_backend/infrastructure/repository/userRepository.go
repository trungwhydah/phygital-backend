package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository struct
type UserRepository struct {
	dbMongo *mongo.Database
}

// NewUserRepository create repository
func NewUserRepository(dbMongo *mongo.Database) *UserRepository {
	return &UserRepository{dbMongo: dbMongo}
}

func (r *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	_, err := r.dbMongo.Collection(user.CollectionName()).InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) CheckExistedEmail(email *string) (bool, error) {
	var user entity.User
	err := u.dbMongo.Collection(user.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "email", Value: *email}}).Decode(&user)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) GetUserByEmail(email *string) (*entity.User, error) {
	var user entity.User
	err := r.dbMongo.Collection(user.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "email", Value: *email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(userID *string) (*entity.User, error) {
	var user entity.User
	err := r.dbMongo.Collection(user.CollectionName()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: *userID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpsertUser(user *entity.User) error {

	var (
		filter     = bson.D{{Key: "_id", Value: user.ID}}
		opts       = options.Update().SetUpsert(true)
		updateData = bson.D{
			{Key: "$set", Value: *user},
		}
	)
	_, err := r.dbMongo.Collection(user.CollectionName()).UpdateOne(
		context.TODO(),
		filter,
		updateData,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateRole(role *string, userID *string) (bool, error) {
	filter := bson.M{"_id": *userID}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: *role}}}}
	result, err := r.dbMongo.Collection(entity.User{}.CollectionName()).UpdateOne(context.Background(), &filter, &update)
	if err != nil {
		return false, err
	}
	return (result.MatchedCount + result.ModifiedCount) != 0, nil
}

func (r *UserRepository) UpdateOrgID(orgID *primitive.ObjectID, userID *string) (bool, error) {
	filter := bson.M{"_id": *userID}
	update := bson.M{"$set": bson.M{"organization_id": *orgID}}
	_, err := r.dbMongo.Collection(entity.User{}.CollectionName()).UpdateOne(context.Background(), &filter, &update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) UpdateUserDetails(userID *string, req *request.UpdateUserDetailsRequest) (bool, error) {
	filter := bson.D{{Key: "_id", Value: *userID}}
	update := bson.D{{Key: "$set", Value: *req}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser entity.User
	err := r.dbMongo.Collection(entity.User{}.CollectionName()).FindOneAndUpdate(context.TODO(), filter, update, option).Decode(&updatedUser)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) GetUserWithNoWallet() (*[]entity.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"wallet_address": bson.M{"$exists": false}},
			{"wallet_address": ""},
			{"wallet_address": nil},
		},
	}
	cursor, err := r.dbMongo.Collection(entity.User{}.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}
