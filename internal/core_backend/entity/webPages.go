package entity

type WebPageBase struct {
	BaseModel `bson:"inline"`
	Name      string `bson:"name,omitempty"`
	URLLink   string `bson:"url_link"`
	Type      string `bson:"type,omitempty"`
}
type WebPage struct {
	WebPageBase `bson:"inline"`
	Attributes  map[string]interface{} `bson:"attributes"`
}

// CollectionName of WebPage
func (WebPage) CollectionName() string {
	return "webpages"
}
