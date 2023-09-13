package scan

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// Scan interface
type Scan interface {
	// Interface for repository
	GetTagWithID(tagpID *string) (*entity.Tag, error)
}

// Repository interface
type Repository interface {
	Scan
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	ProcessScan(*request.ScanRequest) (*entity.Scan, int, error)
}
