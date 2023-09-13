package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateMappingRequest struct {
	ProductItemID  *primitive.ObjectID `json:"product_item_id,omitempty" bson:"product_item_id,omitempty"`
	Claimable      *bool               `json:"claimable,omitempty" bson:"claimable,omitempty"`
	ExternalURL    *string             `json:"external_url,omitempty" bson:"external_url,omitempty"`
	OrgID          *primitive.ObjectID `json:"org_id,omitempty" bson:"org_id,omitempty"`
	DigitalAssetID *primitive.ObjectID `json:"digital_asset_id,omitempty" bson:"digital_asset_id,omitempty"`
}

type MultipleMappingWithSingleProduct struct {
	TagIDList   []string `form:"tag_id_list" validate:"required"`
	ProductID   string   `form:"product_id" validate:"required"`
	ExternalURL string   `form:"external_url" validate:"required"`
}

type UnmapRequest struct {
	ProductItemID string `json:"product_item_id" validate:"required"`
}
