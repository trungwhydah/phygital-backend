package mapping

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mapping interface
type Mapping interface {
	// Interface for repository
	GetAllMapping() (*[]entity.Mapping, error)
	GetAllMappingForProduct(*string, *string) (*[]entity.Mapping, error)
	GetAllMappingInOrg(orgID *primitive.ObjectID) (*[]entity.Mapping, error)
	UpsertMapping(mapping *entity.Mapping) (bool, error)
	GetMappingWithTagID(tagID *string) (*entity.Mapping, error)
	GetMappingWithProductItemID(productItemID *string) (*entity.Mapping, error)
	UpdateMapping(*string, *request.UpdateMappingRequest) (bool, error)
	Unmap(*string) (bool, error)
	GetMappingByDigitalAsset(digitalAssetID *string) (*entity.Mapping, error)
}

// Repository interface
type Repository interface {
	Mapping
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	GetAllMappingInOrg(orgID *string) (*[]entity.Mapping, int, error)
	GetAllMappingForProduct(*string, *string) (*[]entity.Mapping, int, error)
	InitMapping(tagID, orgID *string) (bool, int, error)
	UpdateMapping(*string, *request.UpdateMappingRequest) (bool, int, error)
	Unmap(*string) (bool, int, error)
	GetMappingWithTagID(tagID *string) (*entity.Mapping, int, error)
	IsProductItemIDMapped(productItemID *string) (bool, int, error)
	GetMappingWithProductItemID(productItemID *string) (*entity.Mapping, int, error)
	GetMappingByDigitalAsset(digitalAssetID *string) (*entity.Mapping, int, error)
}
