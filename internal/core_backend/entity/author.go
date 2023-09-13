package entity

type Author struct {
	BaseModel      `bson:",inline" json:",inline"`
	Name           string `bson:"name" json:"name" biding:"required"`
	ExperienceYear string `bson:"experience_year" json:"experience_year" biding:"omitempty"`
	ArtworksCount  string `bson:"artworks_count" json:"artworks_count" biding:"omitempty"`
	Phone          string `bson:"phone" json:"phone"  biding:"required"`
	Email          string `bson:"email" json:"email"  biding:"email"`
	Avatar         Media  `bson:"avatar" json:"avatar"`
	Description    string `bson:"description" json:"description"  biding:"omitempty"`
	Type           string `bson:"type" json:"type" biding:"omitempty"`
}

// CollectionName Collection name of Author
func (Author) CollectionName() string {
	return "authors"
}
