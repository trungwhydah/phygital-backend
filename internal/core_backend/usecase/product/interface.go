package product

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product interface
type Product interface {
	// Interface for repository
	CreateProduct(*entity.Product) (*entity.Product, error)
	CheckExistedProduct(*entity.Product) (bool, error)
	GetAllProducts() (*[]entity.Product, error)
	GetAllProductsInOrg(orgID *primitive.ObjectID) (*[]entity.Product, error)
	GetProductByID(*string) (*entity.Product, error)
	UpdateProductDetail(*entity.Product, *string) (bool, error)
	SoftDeleteProductByID(productID *string) (bool, error)
	UpdateProductTotalItems(*string, int) (bool, error)
	GetProductForAuthor(*string) (*[]entity.Product, error)
}

// Repository interface
type Repository interface {
	Product
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateProduct(*entity.Product) (*entity.Product, int, error)
	GetAllProductsInOrg(orgID *string) (*[]entity.Product, int, error)
	GetProductDetail(*request.InteractProductDetailRequest) (*entity.Product, int, error)
	UpdateProductDetail(*entity.Product, *string) (bool, int, error)
	DeteleProductByID(*request.InteractProductDetailRequest) (bool, int, error)
	GetProductByID(*string) (*entity.Product, int, error)
	SyncTotalItems(*string, int) (bool, int, error)
	GetProductForAuthor(*string) (*[]entity.Product, int, error)
}
