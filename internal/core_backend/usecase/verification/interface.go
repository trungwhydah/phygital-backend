package verification

import (
	"backend-service/internal/core_backend/entity"
)

// verification interface
type Verification interface {
	// Interface for repository
	SaveVerifition(ver *entity.Verification) (*entity.Verification, error)
}

// Repository interface
type Repository interface {
	Verification
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	Verify(scan *entity.Scan, chip *entity.Tag) (*entity.Verification, int, error)
}
