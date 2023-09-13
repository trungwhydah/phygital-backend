package request

type TemplatePagesRequest struct {
	PageID string `json:"page_id"`
}

type TemplateMenuRequest struct {
	Title  TemplateMenuTitleRequest `json:"title"`
	PageID string                   `json:"page_id"`
}

type TemplateMenuTitleRequest struct {
	VI string `json:"vi"`
	EN string `json:"en"`
}

type CreateTemplateRequest struct {
	Name      string                 `json:"name,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Languages []string               `json:"languages"`
	Pages     []TemplatePagesRequest `json:"pages,omitempty"`
	Menu      []TemplateMenuRequest  `json:"menu,omitempty"`
}

type UpdateTemplateRequest struct {
	TemplateID string                 `validate:"required" swaggerignore:"true"`
	Name       string                 `json:"name,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Languages  []string               `json:"languages"`
	Pages      []TemplatePagesRequest `json:"pages,omitempty"`
	Menu       []TemplateMenuRequest  `json:"menu,omitempty"`
}
