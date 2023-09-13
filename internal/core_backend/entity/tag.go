package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	BaseModel      `bson:"inline"`
	HardwareID     string             `bson:"hardware_id"`
	TagID          string             `bson:"tag_id"`
	TagType        string             `bson:"tag_type"` // chip or qr
	EncryptMode    string             `bson:"encrypt_mode"`
	RawData        string             `bson:"raw_data"`
	ScanCounter    int                `bson:"scan_counter"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
}

// CollectionName Collection name of Tag
func (Tag) CollectionName() string {
	return "tags"
}
