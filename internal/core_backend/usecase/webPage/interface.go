package webpage

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// WebPage interface repository methods
type WebPage interface {
	// Interface for repository
	CreateWebPage(page *entity.WebPage) (*entity.WebPage, error)
	GetAllWebPages() (*[]entity.WebPage, error)
	GetWebPage(pageId *string) (*entity.WebPage, error)
	UpdateWebPage(page *entity.WebPage) (*entity.WebPage, error)
	DeleteWebPage(Id *string) (bool, error)
}

// Repository interface
type Repository interface {
	WebPage
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateWebPage(*request.CreateWebpageRequest) (string, int, error)
	GetAllWebPages() (*[]entity.WebPage, int, error)
	GetWebPage(pageId *string) (*entity.WebPage, int, error)
	UpdateWebPage(*request.UpdateWebpageRequest) (string, int, error)
	DeleteWebPage(pageId *string) (bool, int, error)
}
