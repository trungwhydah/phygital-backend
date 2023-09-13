package product

import (
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service struct
type Service struct {
	repo Repository
}

// NewService create service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateProduct
func (s *Service) CreateProduct(request *entity.Product) (*entity.Product, int, error) {
	productInserted, err := s.repo.CreateProduct(request)
	if err != nil {
		logger.LogError("Get error when creating product: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return productInserted, http.StatusOK, nil
}

func (s *Service) GetAllProductsInOrg(orgID *string) (*[]entity.Product, int, error) {
	var products *[]entity.Product
	if len(*orgID) != 0 {
		oID, err := primitive.ObjectIDFromHex(*orgID)
		if err != nil {
			logger.LogError("Got error while parsing organization: " + err.Error())
			return nil, http.StatusInternalServerError, err
		}

		result, err := s.repo.GetAllProductsInOrg(&oID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		products = result
	} else {
		result, err := s.repo.GetAllProducts()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		products = result
	}

	return products, http.StatusOK, nil
}

func (s *Service) GetProductDetail(request *request.InteractProductDetailRequest) (*entity.Product, int, error) {
	product, err := s.repo.GetProductByID(&request.ProductID)
	if err != nil {
		logger.LogError("Got error while get product: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return product, http.StatusOK, nil
}

func (s *Service) UpdateProductDetail(request *entity.Product, productID *string) (bool, int, error) {
	prod, err := s.repo.GetProductByID(productID)
	if err != nil {
		logger.LogError("Got error while getting product detail: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	request.BaseModel = prod.BaseModel
	request.SetTime()
	request.ClearModel()
	success, err := s.repo.UpdateProductDetail(request, productID)
	if err != nil {
		logger.LogError("Got error while updating product detail: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) DeteleProductByID(request *request.InteractProductDetailRequest) (bool, int, error) {
	success, err := s.repo.SoftDeleteProductByID(&request.ProductID)
	if err != nil {
		logger.LogError("Got error while deleting product: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) GetProductByID(productID *string) (*entity.Product, int, error) {
	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		logger.LogError("[DEBUG] - 12 - Error: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return product, http.StatusOK, nil
}

func (s *Service) SyncTotalItems(productID *string, totalItems int) (bool, int, error) {
	ok, err := s.repo.UpdateProductTotalItems(productID, totalItems)
	if err != nil {
		logger.LogError(err.Error())
		return false, http.StatusInternalServerError, err
	}
	return ok, http.StatusOK, nil
}

// GetProductForAuthor -
func (s *Service) GetProductForAuthor(authorID *string) (*[]entity.Product, int, error) {
	listProducts, err := s.repo.GetProductForAuthor(authorID)
	if err != nil {
		logger.LogError("Error getting product: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return listProducts, http.StatusOK, nil
}
