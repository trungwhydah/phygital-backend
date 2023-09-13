package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type DigitalAssetCollection struct {
	BaseModel       `bson:"inline"`
	Name            string             `bson:"name"`
	Chain           string             `bson:"chain"`
	ChainID         int                `bson:"chain_id"`
	Description     string             `bson:"description"`
	ContractAddress string             `bson:"contract_address"`
	Standard        string             `bson:"standard"`
	OrganizationID  primitive.ObjectID `bson:"org_id"`
}

// CollectionName Collection name of DigitalAssetCollection
func (DigitalAssetCollection) CollectionName() string {
	return "digital_asset_collections"
}
