package template

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo Repository
}

// NewService create service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateTemplate(request *request.CreateTemplateRequest) (bool, int, error) {
	var pages []entity.TemplatePages
	for _, page := range request.Pages {
		pageID, err := primitive.ObjectIDFromHex(page.PageID)
		if err != nil {
			logger.LogError("Error convert pageID from string to ObjectID")
			return false, http.StatusInternalServerError, err
		}
		pages = append(pages, entity.TemplatePages{
			PageID: pageID,
		})
	}
	var menus []entity.TemplateMenu
	for _, menu := range request.Menu {
		pageID, err := primitive.ObjectIDFromHex(menu.PageID)
		if err != nil {
			logger.LogError("Error convert pageID from string to ObjectID")
			return false, http.StatusInternalServerError, err
		}
		menus = append(menus, entity.TemplateMenu{
			Title: entity.TemplateMenuTitle{
				VI: menu.Title.VI,
				EN: menu.Title.EN,
			},
			PageID: pageID,
		})
	}
	template := &entity.Template{
		Name:      request.Name,
		Category:  request.Category,
		Languages: request.Languages,
		Pages:     pages,
		Menu:      menus,
	}
	template.SetTime()
	templateInserted, err := s.repo.CreateTemplate(template)
	if err != nil {
		logger.LogError("Error creating template: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	return !templateInserted.ID.IsZero(), http.StatusOK, nil
}

func (s *Service) UpdateTemplate(request *request.UpdateTemplateRequest) (bool, int, error) {
	template, err := s.repo.GetTemplate(&request.TemplateID)
	if err != nil {
		logger.LogError("Error getting template: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	var pages []entity.TemplatePages
	for _, page := range request.Pages {
		pageID, err := primitive.ObjectIDFromHex(page.PageID)
		if err != nil {
			logger.LogError("Error convert pageID from string to ObjectID")
			return false, http.StatusInternalServerError, err
		}
		pages = append(pages, entity.TemplatePages{
			PageID: pageID,
		})
	}
	var menus []entity.TemplateMenu
	for _, menu := range request.Menu {
		pageID, err := primitive.ObjectIDFromHex(menu.PageID)
		if err != nil {
			logger.LogError("Error convert pageID from string to ObjectID")
			return false, http.StatusInternalServerError, err
		}
		menus = append(menus, entity.TemplateMenu{
			Title: entity.TemplateMenuTitle{
				VI: menu.Title.VI,
				EN: menu.Title.EN,
			},
			PageID: pageID,
		})
	}

	template.Name = request.Name
	template.Category = request.Category
	template.Languages = request.Languages
	template.Pages = pages
	template.Menu = menus
	template.SetTime()
	err = s.repo.UpdateTemplate(template)
	if err != nil {
		logger.LogError("Error updating template: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	return true, http.StatusOK, nil
}

func (s *Service) GetTemplate(templateID *string) (*entity.Template, int, error) {
	template, err := s.repo.GetTemplate(templateID)
	if err != nil {
		logger.LogError("[DEBUG] - 14 - Error getting template: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	return template, http.StatusOK, nil
}

func (s *Service) CheckExistedTemplate(templateID *string) (bool, int, error) {
	isExisted, err := s.repo.CheckExistedTemplate(templateID)
	if err != nil {
		logger.LogError("Error when checking template existence: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	return isExisted, http.StatusOK, err
}

func (s *Service) GetAllTemplates() (*[]entity.Template, int, error) {
	templates, err := s.repo.GetAllTemplates()
	if err != nil {
		logger.LogError("Error getting all templates: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	return templates, http.StatusOK, nil
}

func (s *Service) GetTemplateWebpages(tID *string) (*entity.TemplateWebpages, int, error) {
	templateWebpages, err := s.repo.GetTemplateWebpages(tID)
	if err != nil {
		logger.LogError("Error getting template and its webpages: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	return templateWebpages, http.StatusOK, nil
}

// CloneTemplate - Clone a template with given template ID
func (s *Service) CloneTemplate(templateID *string) (*entity.Template, int, error) {
	template, err := s.repo.GetTemplate(templateID)
	if err != nil {
		logger.LogError("Error getting template: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	template.Renew()
	template.Name += "_Copy"

	templateInserted, err := s.repo.CreateTemplate(template)
	if err != nil {
		logger.LogError("Error creating template: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	return templateInserted, http.StatusOK, nil
}
