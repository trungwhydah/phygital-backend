package template

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// Repository Interface For Template
type Template interface {
	CreateTemplate(*entity.Template) (*entity.Template, error)
	UpdateTemplate(*entity.Template) error
	GetTemplate(*string) (*entity.Template, error)
	CheckExistedTemplate(*string) (bool, error)
	GetAllTemplates() (*[]entity.Template, error)
	GetTemplateWebpages(tID *string) (*entity.TemplateWebpages, error)
}

type Repository interface {
	Template
}

type Usecase interface {
	CreateTemplate(*request.CreateTemplateRequest) (bool, int, error)
	UpdateTemplate(*request.UpdateTemplateRequest) (bool, int, error)
	GetTemplate(*string) (*entity.Template, int, error)
	CheckExistedTemplate(*string) (bool, int, error)
	GetAllTemplates() (*[]entity.Template, int, error)
	GetTemplateWebpages(tID *string) (*entity.TemplateWebpages, int, error)
}
