package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DigitalAsset struct {
	BaseModel    `bson:"inline"`
	CollectionID primitive.ObjectID `bson:"collection_id"`
	TokenID      int64                `bson:"token_id"`
	TxHash       string             `bson:"tx_hash"`
	OwnerAddress string             `bson:"owner_address"`
	Metadata     Metadata           `bson:"metadata"`
}

// CollectionName Collection name of DigitalAsset
func (DigitalAsset) CollectionName() string {
	return "digital_assets"
}

type DigitalAssetProductAggregate struct {
	DigitalAsset `bson:"inline"`
	ItemIndex    int     `bson:"item_index"`
	OrgTagName   string  `bson:"org_tag_name"`
	Product      Product `bson:"product"`
}

type Metadata struct {
	Description  string              `json:"description" bson:"description"`
	ExternalURL  string              `json:"external_url" bson:"external_url"`
	Image        string              `json:"image" bson:"image"`
	Name         string              `json:"name" bson:"name"`
	Attributes   []MetadataAttribute `json:"attributes" bson:"attributes"`
	AnimationURL string              `json:"animation_url" bson:"animation_url"`
}

type MetadataAttribute struct {
	TraitType   string `json:"trait_type" bson:"trait_type"`
	Value       string `json:"value" bson:"value"`
	DisplayType string `json:"display_type" bson:"display_type"`
}
