package upload

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	UploadImagesToGCPStorage(*request.UploadImagesRequest) (*[]entity.Upload, int, error)
}
