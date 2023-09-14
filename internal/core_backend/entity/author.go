package entity

type Author struct {
	BaseModel      `bson:",inline" json:",inline"`
	Name           MultipleLanguages `bson:"name" json:"name" binding:"required"`
	ExperienceYear string            `bson:"experience_year" json:"experience_year" binding:"omitempty"`
	ArtworksCount  string            `bson:"artworks_count" json:"artworks_count" binding:"omitempty"`
	Phone          string            `bson:"phone" json:"phone"  binding:"omitempty"`
	Email          string            `bson:"email" json:"email"  binding:"omitempty,email"`
	Avatar         Media             `bson:"avatar" json:"avatar"`
	Translation    map[string]any    `bson:"translation" json:"translation"`
	Type           string            `bson:"type" json:"type" binding:"omitempty"`
	ContactName    string            `bson:"contact_name" json:"contact_name"`
}

// CollectionName Collection name of Author
func (Author) CollectionName() string {
	return "authors"
}
