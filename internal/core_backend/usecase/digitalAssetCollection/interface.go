package digitalAssetCollection

import (
	"backend-service/internal/core_backend/entity"
)

// DigitalAssetCollection interface
type DigitalAssetCollection interface {
	// Interface for repository
	GetCollectionByOrgID(orgID *string) (*entity.DigitalAssetCollection, error)
	GetCollectionByID(cID *string) (*entity.DigitalAssetCollection, error)
}

// Repository interface
type Repository interface {
	DigitalAssetCollection
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	GetCollectionByOrgID(orgID *string) (*entity.DigitalAssetCollection, int, error)
	GetCollectionByID(cID *string) (*entity.DigitalAssetCollection, int, error)
}
