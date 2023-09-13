package digitalAsset

import (
	"backend-service/internal/core_backend/entity"
)

// DigitalAsset interface
type DigitalAsset interface {
	// Interface for repository
	GetDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, error)
	CreateDigitalAsset(da *entity.DigitalAsset) (*entity.DigitalAsset, error)
	UpdateDigitalAsset(da *entity.DigitalAsset) (bool, error)
	UpdateDigitalAssetMetadata(*string, *entity.Metadata) (bool, error)
	GetDigitalAssetByID(daID *string) (*entity.DigitalAsset, error)
	GetDigitalAssetByTokenID(collectionID *string, tokenID *int) (*entity.DigitalAsset, error)
	GetAllActiveDigitalAssets() (*[]entity.DigitalAsset, error)
	GetActiveDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, error)
	GetDigitalAssetsProductAggregate() (*[]entity.DigitalAssetProductAggregate, error)
}

// Repository interface
type Repository interface {
	DigitalAsset
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	GetDigitalAssetByCollection(collectionID *string) (*[]entity.DigitalAsset, int, error)
	CreateDigitalAsset(da *entity.DigitalAsset) (*entity.DigitalAsset, int, error)
	UpdateDigitalAsset(da *entity.DigitalAsset) (bool, int, error)
	UpdateDigitalAssetMetadata(*string, *entity.Metadata) (bool, int, error)
	GetDigitalAssetByID(daID *string) (*entity.DigitalAsset, int, error)
	GetDigitalAssetByTokenID(collectionID *string, tokenID *int) (*entity.DigitalAsset, int, error)
	GetAllActiveDigitalAssets() (*[]entity.DigitalAsset, int, error)
	GetActiveDigitalAssetByCollectionID(collectionID *string) (*[]entity.DigitalAsset, int, error)
	GetDigitalAssetsProductAggregate() (*[]entity.DigitalAssetProductAggregate, int, error)
	ConstructMetadata(int, *string, *entity.Product) *entity.Metadata
}
