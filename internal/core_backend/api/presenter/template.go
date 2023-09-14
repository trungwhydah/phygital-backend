package presenter

import (
	"backend-service/internal/core_backend/entity"
	"time"
)

type TemplateWebpagesResponse struct {
	ID        string                          `json:"template_id"`
	Name      string                          `json:"name"`
	Category  string                          `json:"category"`
	Languages []string                        `json:"languages"`
	Pages     []TemplateWebpagesPagesResponse `json:"pages"`
	Menu      []TemplateWebpagesMenuResponse  `json:"menu"`
}

type TemplateWebpagesPagesResponse struct {
	ID         string                 `json:"page_id"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Name       string                 `json:"name"`
	URLLink    string                 `json:"url_link"`
	Type       string                 `json:"type"`
	Category   string                 `json:"category"`
	Attributes map[string]interface{} `json:"attributes"`
}

type TemplateWebpagesMenuResponse struct {
	Title      TemplateMenuTitleResponse `json:"title"`
	ID         string                    `json:"page_id"`
	Status     string                    `json:"status"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
	Name       string                    `json:"name"`
	URLLink    string                    `json:"url_link"`
	Type       string                    `json:"type"`
	Category   string                    `json:"category"`
	Attributes map[string]interface{}    `json:"attributes"`
}

type TemplateResponse struct {
	ID        string                  `json:"template_id"`
	Name      string                  `json:"name"`
	Category  string                  `json:"category"`
	Languages []string                `json:"languages"`
	Pages     []TemplatePagesResponse `json:"pages"`
	Menu      []TemplateMenuResponse  `json:"menu"`
}

type TemplatePagesResponse struct {
	PageID string `json:"page_id"`
}

type TemplateMenuResponse struct {
	Title  TemplateMenuTitleResponse `json:"title"`
	PageID string                    `json:"page_id"`
}

type TemplateMenuTitleResponse struct {
	VI string `json:"vi"`
	EN string `json:"en"`
}

type ListTemplateResponse struct {
	TemplateList []TemplateResponse `json:"template_list"`
}

type PresenterTemplate struct{}

type ConvertTemplate interface {
	ResponseGetTemplate(template *entity.Template) *TemplateResponse
	ResponseGetTemplateWebpages(templateWebpages *entity.TemplateWebpages) *TemplateWebpagesResponse
	ResponseGetAllTemplates(templates *[]entity.Template) *ListTemplateResponse
}

func NewPresenterTemplate() ConvertTemplate {
	return &PresenterTemplate{}
}

func (pt *PresenterTemplate) ResponseGetTemplate(template *entity.Template) *TemplateResponse {
	var pages []TemplatePagesResponse
	for _, page := range template.Pages {
		pages = append(pages, TemplatePagesResponse{
			PageID: page.PageID.Hex(),
		})
	}
	var menus []TemplateMenuResponse
	for _, menu := range template.Menu {
		menus = append(menus, TemplateMenuResponse{
			Title: TemplateMenuTitleResponse{
				VI: menu.Title.VI,
				EN: menu.Title.EN,
			},
			PageID: menu.PageID.Hex(),
		})
	}
	return &TemplateResponse{
		ID:        template.ID.Hex(),
		Name:      template.Name,
		Category:  template.Category,
		Languages: template.Languages,
		Pages:     pages,
		Menu:      menus,
	}
}

func (pt *PresenterTemplate) ResponseGetTemplateWebpages(templateWebpages *entity.TemplateWebpages) *TemplateWebpagesResponse {
	var pages []TemplateWebpagesPagesResponse
	for _, page := range templateWebpages.Pages {
		pages = append(pages, TemplateWebpagesPagesResponse{
			ID:         page.ID.Hex(),
			Status:     page.Status,
			CreatedAt:  page.CreatedAt,
			UpdatedAt:  page.UpdatedAt,
			Name:       page.Name,
			URLLink:    page.URLLink,
			Type:       page.Type,
			Attributes: page.Attributes,
		})
	}
	var menus []TemplateWebpagesMenuResponse
	for _, menu := range templateWebpages.Menu {
		menus = append(menus, TemplateWebpagesMenuResponse{
			Title: TemplateMenuTitleResponse{
				VI: menu.Title.VI,
				EN: menu.Title.EN,
			},
			ID:         menu.ID.Hex(),
			Status:     menu.Status,
			CreatedAt:  menu.CreatedAt,
			UpdatedAt:  menu.UpdatedAt,
			Name:       menu.Name,
			URLLink:    menu.URLLink,
			Type:       menu.Type,
			Attributes: menu.Attributes,
		})
	}
	return &TemplateWebpagesResponse{
		ID:        templateWebpages.ID.Hex(),
		Name:      templateWebpages.Name,
		Category:  templateWebpages.Category,
		Languages: templateWebpages.Languages,
		Pages:     pages,
		Menu:      menus,
	}

}

func (pt *PresenterTemplate) ResponseGetAllTemplates(templates *[]entity.Template) *ListTemplateResponse {
	var response ListTemplateResponse
	for _, template := range *templates {
		var pages []TemplatePagesResponse
		for _, page := range template.Pages {
			pages = append(pages, TemplatePagesResponse{
				PageID: page.PageID.Hex(),
			})
		}
		var menus []TemplateMenuResponse
		for _, menu := range template.Menu {
			menus = append(menus, TemplateMenuResponse{
				Title: TemplateMenuTitleResponse{
					VI: menu.Title.VI,
					EN: menu.Title.EN,
				},
				PageID: menu.PageID.Hex(),
			})
		}
		response.TemplateList = append(response.TemplateList, TemplateResponse{
			ID:        template.ID.Hex(),
			Name:      template.Name,
			Category:  template.Category,
			Languages: template.Languages,
			Pages:     pages,
			Menu:      menus,
		})
	}
	return &response
}
