package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var AllowType = []string{"coffee", "astronaut", "sculpture", "ortho"}

type Product struct {
	BaseModel      `bson:"inline"`
	Type           string             `bson:"type" json:"type" binding:"required"`
	ProductName    string             `bson:"product_name" json:"product_name"`
	Origin         string             `bson:"origin" json:"origin"`
	URLLink        string             `bson:"url_link" json:"url_link"`
	TotalItem      int                `bson:"total_item" json:"total_item"`
	TemplateID     primitive.ObjectID `bson:"template_id" json:"template_id"`
	OrganizationID primitive.ObjectID `bson:"org_id" json:"org_id"`
	RatingScore    float64            `bson:"rating_score" json:"rating_score"`
	Image          Media              `bson:"image" json:"image"`
	Video          Media              `bson:"video" json:"video"`
	ThreeDimension Media              `bson:"three_dimension" json:"three_dimension"`
	Tags           []string           `bson:"tags" json:"tags"`
	AuthorID       primitive.ObjectID `bson:"author_id" json:"author_id"`
	Attribute      any                `bson:"attribute" json:"attribute"`
}

// CollectionName Collection name of Product
func (Product) CollectionName() string {
	return "products"
}

type AttributeCoffee struct {
	FarmName       string         `bson:"farm_name" json:"farm_name"`
	FarmVideo      Media          `bson:"farm_video" json:"farm_video"`
	FarmImage      Media          `bson:"farm_image" json:"farm_image"`
	FarmHeight     string         `bson:"farm_height" json:"farm_height"`
	FarmArea       string         `bson:"farm_area" json:"farm_area"`
	Varietal       string         `bson:"varietal" json:"varietal"`
	Process        string         `bson:"process" json:"process"`
	CountryName    string         `bson:"country_name" json:"country_name"`
	CountryVideo   Media          `bson:"country_video" json:"country_video"`
	CountryImage   Media          `bson:"country_image" json:"country_image"`
	BrewingTime    string         `bson:"brewing_time" json:"brewing_time"`
	LinkBuyProduct string         `bson:"link_buy_product" json:"link_buy_product"`
	Acidity        string         `bson:"acidity" json:"acidity"`
	Bitter         string         `bson:"bitter" json:"bitter"`
	Sweet          string         `bson:"sweet" json:"sweet"`
	Translation    map[string]any `bson:"translation" json:"translation"`
}

type AttributeAstronaut struct {
}

type AttributeSculpture struct {
	ContactName           string `bson:"contact_name" json:"contact_name"`
	SculptureRank         string `bson:"sculpture_rank" json:"sculpture_rank"`
	SculpturePedestalSize string `bson:"sculpture_pedestal_size" json:"sculpture_pedestal_size"`
	SculptureSize         string `bson:"sculpture_size" json:"sculpture_size"`
	SculptureTime         string `bson:"sculpture_time" json:"sculpture_time"`
	SculptureWeight       string `bson:"sculpture_weight" json:"sculpture_weight"`
	// remove 3 lines below
	SculptureHeight string `bson:"sculpture_height" json:"sculpture_height"`
	SculptureLength string `bson:"sculpture_length" json:"sculpture_length"`
	SculptureWidth  string `bson:"sculpture_width" json:"sculpture_width"`
	//Description     string         `bson:"description" json:"description"`
	Village     Village        `bson:"village" json:"village"`
	Craftsman   Craftsman      `bson:"craftsman" json:"craftsman"`
	Stone       Stone          `bson:"stone" json:"stone"`
	Processes   []Process      `bson:"processes" json:"processes"`
	Translation map[string]any `bson:"translation" json:"translation"`
}

type AttributeOrtho struct {
}

type Craftsman struct {
	Name           string `bson:"name" json:"name"`
	ExperienceYear string `bson:"experience_year" json:"experience_year"`
	ArtworksCount  string `bson:"artworks_count" json:"artworks_count"`
	Phone          string `bson:"phone" json:"phone"`
	Email          string `bson:"email" json:"email"`
	Avatar         Media  `bson:"avatar" json:"avatar"`
	Description    string `bson:"description" json:"description"`
}

type Process struct {
	ImageURL    string `bson:"image_url" json:"image_url"`
	Description struct {
		EN string `bson:"en" json:"en"`
		VI string `bson:"vi" json:"vi"`
	} `bson:"description" json:"description"`
}

type Stone struct {
	//Name       string `bson:"name" json:"name"`
	Origin string `bson:"origin" json:"origin"`
	//Clarity    string `bson:"clarity" json:"clarity"`
	//Rarity     string `bson:"rarity" json:"rarity"`
	//Properties string `bson:"properties" json:"properties"`
	//Color      string `bson:"color" json:"color"`
	Image       Media          `bson:"image" json:"image"`
	Translation map[string]any `bson:"translation" json:"translation"`
}

type Village struct {
	Name          string         `bson:"name" json:"name"`
	Translation   map[string]any `bson:"translation" json:"translation"`
	LocationVideo Media          `bson:"location_video" json:"location_video"`
}

// ParseAttribute
func (p *Product) ParseAttribute() *Product {
	switch p.Type {
	case "coffee":
		marAtt, err := bson.Marshal(p.Attribute)
		if err != nil {
			return p
		}

		var coffeeAttribute AttributeCoffee
		bson.Unmarshal(marAtt, &coffeeAttribute)
		p.Attribute = coffeeAttribute
	case "astronaut":
		p.Attribute = AttributeAstronaut{}
	case "sculpture":
		marAtt, err := bson.Marshal(p.Attribute)
		if err != nil {
			return p
		}

		var sculptureAttribute AttributeSculpture
		bson.Unmarshal(marAtt, &sculptureAttribute)
		p.Attribute = sculptureAttribute
	case "ortho":
		p.Attribute = AttributeOrtho{}
	}

	return p
}
