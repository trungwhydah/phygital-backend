package webpage

import (
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"
)

// Service struct
type Service struct {
	repo Repository
}

// NewService create service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateWebpage
func (s *Service) CreateWebPage(request *request.CreateWebpageRequest) (string, int, error) {
	page := &entity.WebPage{
		WebPageBase: entity.WebPageBase{
			Name:    request.Name,
			URLLink: request.URLLink,
			Type:    request.Type,
		},
		Attributes: request.Attributes,
	}
	webPage, err := s.repo.CreateWebPage(page)
	if err != nil {
		logger.LogError("Get error when getting webpage: " + err.Error())
		return "", http.StatusInternalServerError, err
	}

	return webPage.ID.Hex(), http.StatusOK, nil
}

// GetWebPage
func (s *Service) GetWebPage(pageId *string) (*entity.WebPage, int, error) {

	webPage, err := s.repo.GetWebPage(pageId)
	if err != nil {
		logger.LogError("Get error when getting webpage: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return webPage, http.StatusOK, nil
}

// UpdateWebpage
func (s *Service) UpdateWebPage(request *request.UpdateWebpageRequest) (string, int, error) {
	page, err := s.repo.GetWebPage(&request.WebpageID)
	if err != nil {
		logger.LogError("Error getting webpage: " + err.Error())
		return "", http.StatusInternalServerError, err
	}
	page.Name = request.Name
	page.URLLink = request.URLLink
	page.Type = request.Type
	page.Attributes = request.Attributes
	webPage, err := s.repo.UpdateWebPage(page)
	if err != nil {
		logger.LogError("Get error when updating webpage: " + err.Error())
		return "", http.StatusInternalServerError, err
	}

	return webPage.ID.Hex(), http.StatusOK, nil
}

// DeleteWebPage
func (s *Service) DeleteWebPage(pageId *string) (bool, int, error) {

	result, err := s.repo.DeleteWebPage(pageId)
	if err != nil {
		logger.LogError("Get error when deleting webpage: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}

// GetWebPage
func (s *Service) GetAllWebPages() (*[]entity.WebPage, int, error) {

	webPages, err := s.repo.GetAllWebPages()
	if err != nil {
		logger.LogError("Get error when getting all webpages: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return webPages, http.StatusOK, nil
}
