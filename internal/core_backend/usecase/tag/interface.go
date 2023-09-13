package tag

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// Tag interface
type Tag interface {
	// Interface for repository
	CreateTag(*entity.Tag) (*entity.Tag, error)
	CheckExistedTag(*entity.Tag) (bool, error)
	GetTagNotMapped(tagMapped *[]string) (*[]entity.Tag, error)
	GetTag(string) (*entity.Tag, error)
	GetTagByHWID(*string) (*entity.Tag, error)
	UpdateTagCounter(tagID *string, scanCounter *int) (bool, error)
}

// Repository interface
type Repository interface {
	Tag
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateTag(*request.CreateTagRequest) (bool, int, error)
	GetTagNotMapped(tagMapped *[]string) (*[]entity.Tag, int, error)
	GetTag(string) (*entity.Tag, int, error)
	UpdateTagCounter(tagID *string, scanCounter *int) (bool, int, error)
	GetTagByHWID(uid *string) (*entity.Tag, int, error)
}
