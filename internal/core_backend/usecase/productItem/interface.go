package productItem

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductItem interface
type ProductItem interface {
	// Interface for repository
	InsertProductItem(*entity.ProductItem) (*entity.ProductItem, error)
	CheckProductItemMapped(productItemID *string) (bool, error)
	GetAllProductItemByProductID(productID *primitive.ObjectID) (*[]entity.ProductItem, error)
	GetAllProductItem() (*[]entity.ProductItem, error)
	GetDetailProductItemByID(*string) (*entity.ProductItem, error)
	IsAbleToClaim(productItemID *string) (bool, error)
	SetOwner(productItemID, ownerID *string) (bool, error)
	ToggleClaimable(productItemID *string) (bool, error)
	GetDetailWithTagID(tagID *string) (*entity.ProductItem, error)
	UpdateTotalLike(productItemID *string) (bool, error)
	GetOrganizationNameByProductItemID(productItemID *string) (string, error)
	CountNumProductItems(*string) (int, error)
	GetProductItemsInOrg(*string) (*[]entity.ProductItem, error)
	GetProductItemProductOrgAggregate(*string) (*entity.ProductItemProductOrgAggregate, error)
}

// Repository interface
type Repository interface {
	ProductItem
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateProductItem(*request.CreateProductItemRequest) (*entity.ProductItem, int, error)
	GetAllProductItem() (*[]entity.ProductItem, int, error)
	GetAllProductItemInProduct(productID *string) (*[]entity.ProductItem, int, error)
	GetDetailProductItem(*string) (*entity.ProductItem, int, error)
	SetOnwerForItem(*request.SetOwnerRequest) (bool, int, error)
	ToggleClaimable(*request.ProductItemInteractionRequest) (bool, int, error)
	CheckProductItemMapped(productItemID *string) (bool, int, error)
	UpdateTotalLike(productItemID *string) (bool, int, error)
	CreateMultipleProductItems(*string, int, int) (bool, int, error)
	CountNumProductItems(*string) (int, int, error)
	GetProductItemsInOrg(*string) (*[]entity.ProductItem, int, error)
	GetProductItemProductOrgAggregate(*string) (*entity.ProductItemProductOrgAggregate, int, error)
}
