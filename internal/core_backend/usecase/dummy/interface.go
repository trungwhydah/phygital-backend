package dummy

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// Dummy interface
type Dummy interface {
	// Interface for repository
	GetDummyByID(DummyID *int) (*entity.Dummy, error)
}

// Repository interface
type Repository interface {
	Dummy
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	GetDummy(dummyRquest *request.DummyRequest) (*entity.Dummy, int, error)
}
