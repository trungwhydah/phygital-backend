package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mapping struct {
	BaseModel      `bson:"inline"`
	ProductItemID  primitive.ObjectID `bson:"product_item_id"`
	TagID          string             `bson:"tag_id"`
	ExternalURL    string             `bson:"external_url"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	OwnerID        string             `bson:"owner_id"`
	Claimable      bool               `bson:"claimable"`
	DigitalAssetID primitive.ObjectID `bson:"digital_asset_id"`
	IsMinted       bool               `bson:"is_minted"`
}

// CollectionName Collection name of Mapping
func (Mapping) CollectionName() string {
	return "mappings"
}
