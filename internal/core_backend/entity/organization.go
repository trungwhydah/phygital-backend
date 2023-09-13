package entity

type Organization struct {
	BaseModel        `bson:"inline"`
	OrganizationName string `bson:"org_name"`
	NameTag          string `bson:"org_tag_name"`
	LogoURL          string `bson:"org_logo_url"`
	OwnerID          string `bson:"owner_id"`
}

// CollectionName Collection name of Organization
func (Organization) CollectionName() string {
	return "organizations"
}
