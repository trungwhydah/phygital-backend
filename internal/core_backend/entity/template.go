package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	BaseModel       `bson:"inline"`
	Name            string          `bson:"name"`
	Category        string          `bson:"category"`
	CreatedByUserID string          `bson:"created_by_user_id"`
	Languages       []string        `bson:"languages"`
	Pages           []TemplatePages `bson:"pages"`
	Menu            []TemplateMenu  `bson:"menu"`
}

type TemplatePages struct {
	PageID primitive.ObjectID `bson:"page_id"`
}

type TemplateMenu struct {
	Title  TemplateMenuTitle  `bson:"title"`
	PageID primitive.ObjectID `bson:"page_id"`
}

type TemplateMenuTitle struct {
	VI string `bson:"vi"`
	EN string `bson:"en"`
}

type TemplateWebpages struct {
	BaseModel       `bson:"inline"`
	Name            string                 `bson:"name"`
	Category        string                 `bson:"category"`
	CreatedByUserID string                 `bson:"created_by_user_id"`
	Languages       []string               `bson:"languages"`
	Pages           []WebPage              `bson:"pages"`
	Menu            []TemplateWebpagesMenu `bson:"menu"`
}

type TemplateWebpagesMenu struct {
	Title   TemplateMenuTitle `bson:"title"`
	WebPage `bson:"inline"`
}

// CollectionName Collection name of Template
func (Template) CollectionName() string {
	return "templates"
}
