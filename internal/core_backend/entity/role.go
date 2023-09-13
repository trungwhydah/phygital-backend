package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	BaseModel      `bson:"inline"`
	OrganizationID primitive.ObjectID `bson:"org_id"`
	RoleName       string             `bson:"role_name"`
}

// CollectionName Collection name of Role
func (Role) CollectionName() string {
	return "roles"
}
