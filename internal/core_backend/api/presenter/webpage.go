package presenter

import (
	"backend-service/internal/core_backend/entity"
	"time"
)

type WebpageDetailResponse struct {
	ID         string                 `json:"id"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Name       string                 `json:"name"`
	URLLink    string                 `json:"url_link"`
	Type       string                 `json:"type"`
	Category   string                 `json:"category"`
	Attributes map[string]interface{} `json:"attributes"`
}

type AllWebpagesResponse struct {
	WebpagesList []WebpageDetailResponse `json:"webpages_list"`
}

// presenterScan struct
type PresenterWebpage struct{}

// presenterScan interface
type ConvertWebpage interface {
	ResponseWebpageDetail(webpage *entity.WebPage) WebpageDetailResponse
	ResponseAllWebpages(webpages *[]entity.WebPage) AllWebpagesResponse
}

// NewPresenterScan Constructs presenter
func NewPresenterWebpage() ConvertWebpage {
	return &PresenterWebpage{}
}

// Return property data response
func (pw *PresenterWebpage) ResponseWebpageDetail(webpage *entity.WebPage) WebpageDetailResponse {
	return WebpageDetailResponse{
		ID:         webpage.ID.Hex(),
		Status:     webpage.Status,
		CreatedAt:  webpage.CreatedAt,
		UpdatedAt:  webpage.UpdatedAt,
		Name:       webpage.Name,
		URLLink:    webpage.URLLink,
		Type:       webpage.Type,
		Attributes: webpage.Attributes,
	}
}

func (pw *PresenterWebpage) ResponseAllWebpages(webpages *[]entity.WebPage) AllWebpagesResponse {
	var response AllWebpagesResponse
	for _, webpage := range *webpages {
		response.WebpagesList = append(response.WebpagesList, WebpageDetailResponse{
			ID:         webpage.ID.Hex(),
			Status:     webpage.Status,
			CreatedAt:  webpage.CreatedAt,
			UpdatedAt:  webpage.UpdatedAt,
			Name:       webpage.Name,
			URLLink:    webpage.URLLink,
			Type:       webpage.Type,
			Attributes: webpage.Attributes,
		})
	}
	return response
}
