package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductItem struct {
	BaseModel `bson:"inline"`
	ProductID primitive.ObjectID `bson:"product_id"`
	OwnerID   string             `bson:"owner_id"`
	TotalLike int                `bson:"total_like"`
	ItemIndex int                `bson:"item_index"`
}

// CollectionName Collection name of ProductItem
func (ProductItem) CollectionName() string {
	return "product_items"
}

type ProductItemProductOrgAggregate struct {
	ProductItem `bson:"inline"`
	OrgTagName  string  `bson:"org_tag_name"`
	Product     Product `bson:"product"`
}
