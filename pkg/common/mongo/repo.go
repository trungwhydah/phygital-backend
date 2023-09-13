package mongo

import "go.mongodb.org/mongo-driver/mongo"

type RepoInterface interface {
	GetCollection() *mongo.Collection
}
