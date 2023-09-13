package product

import (
	cmentity "backend-service/internal/common/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	cmentity.Entity `bson:"inline"`
	Type            string             `bson:"type" json:"type"`
	ProductName     string             `bson:"product_name" json:"product_name"`
	Origin          string             `bson:"origin" json:"origin"`
	URLLink         string             `bson:"url_link" json:"url_link"`
	TotalItem       int                `bson:"total_item" json:"total_item"`
	TemplateID      primitive.ObjectID `bson:"template_id" json:"template_id"`
	OrganizationID  primitive.ObjectID `bson:"org_id" json:"org_id"`
	RatingScore     float64            `bson:"rating_score" json:"rating_score"`
	Image           cmentity.Media     `bson:"image" json:"image"`
	Video           cmentity.Media     `bson:"video" json:"video"`
	ThreeDimension  cmentity.Media     `bson:"three_dimension" json:"three_dimension"`
	Tags            []string           `bson:"tags" json:"tags"`
	AuthorID        primitive.ObjectID `bson:"author_id" json:"author_id"`
	Attribute       any                `bson:"attribute" json:"attribute"`
}
