package request

type CreateWebpageRequest struct {
	Name       string                 `json:"name,omitempty"`
	URLLink    string                 `json:"url_link"`
	Type       string                 `json:"type"`
	Category   string                 `json:"category"`
	Attributes map[string]interface{} `json:"attributes"`
}

type UpdateWebpageRequest struct {
	WebpageID  string                 `validate:"required" swaggerignore:"true"`
	Name       string                 `json:"name,omitempty"`
	URLLink    string                 `json:"url_link,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
