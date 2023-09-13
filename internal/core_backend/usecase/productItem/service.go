package productItem

import (
	"errors"
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common"
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

// GetItem get productItem info
func (s *Service) CreateProductItem(itemRquest *request.CreateProductItemRequest) (*entity.ProductItem, int, error) {
	productID, err := primitive.ObjectIDFromHex(itemRquest.ProductID)
	if err != nil {
		logger.LogError("Error convert product id string to ObjectID")
		return nil, http.StatusInternalServerError, err
	}
	item := &entity.ProductItem{
		ProductID: productID,
		BaseModel: entity.BaseModel{
			Status: common.StatusActive,
		},
	}
	item.SetTime()

	itemInserted, err := s.repo.InsertProductItem(item)
	if err != nil {
		logger.LogError("Get error when inserting productItem" + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return itemInserted, http.StatusOK, nil
}

// GetAllProductItem
func (s *Service) GetAllProductItem() (*[]entity.ProductItem, int, error) {
	productItems, err := s.repo.GetAllProductItem()
	if err != nil {
		logger.LogError("Get error when getting all product items" + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return productItems, http.StatusOK, nil
}

// GetAllProductItemInAProductSeries
func (s *Service) GetAllProductItemInProduct(productID *string) (*[]entity.ProductItem, int, error) {
	var productItemList *[]entity.ProductItem
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		logger.LogError("Got error while parsing product: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	productItemList, err = s.repo.GetAllProductItemByProductID(&pID)
	if err != nil {
		logger.LogError("Got error while getting product items: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	return productItemList, http.StatusOK, nil
}

// GetDetailProductItem
func (s *Service) GetDetailProductItem(itemID *string) (*entity.ProductItem, int, error) {
	item, err := s.repo.GetDetailProductItemByID(itemID)
	if err != nil {
		logger.LogError("[DEBUG] - 11 - Got error when getting detail product item: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return item, http.StatusOK, nil
}

func (s *Service) SetOnwerForItem(request *request.SetOwnerRequest) (bool, int, error) {
	able, err := s.repo.IsAbleToClaim(&request.ProductItemID)

	if err != nil {
		logger.LogError("got error when setting onwer for item: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	if !able {
		logger.LogInfo(common.MessageErrorNotAbleToClaim)
		return false, http.StatusBadRequest, errors.New(common.MessageErrorNotAbleToClaim)
	}

	success, err := s.repo.SetOwner(&request.ProductItemID, &request.OwnerID)
	if err != nil {
		logger.LogError("got error when setting onwer for item: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) GetDetailWithTagID(tagID *string) (*entity.ProductItem, int, error) {
	item, err := s.repo.GetDetailWithTagID(tagID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return item, http.StatusOK, nil
}

func (s *Service) CheckProductItemMapped(productItemID *string) (bool, int, error) {
	isMapped, err := s.repo.CheckProductItemMapped(productItemID)
	if err != nil {
		logger.LogError("Error when checking product item mapped or not: " + err.Error())
		return false, http.StatusBadRequest, err
	}
	return isMapped, http.StatusOK, err
}

func (s *Service) UpdateTotalLike(productItemID *string) (bool, int, error) {
	ok, err := s.repo.UpdateTotalLike(productItemID)
	if err != nil {
		return ok, http.StatusInternalServerError, err
	}

	return ok, http.StatusOK, nil
}

func (s *Service) ToggleClaimable(req *request.ProductItemInteractionRequest) (bool, int, error) {
	success, err := s.repo.ToggleClaimable(&req.ProductItemID)
	if err != nil {
		logger.LogError("Got error when setting claimable for item: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) CreateMultipleProductItems(productID *string, totalItems int, startIndex int) (bool, int, error) {
	pID, err := primitive.ObjectIDFromHex(*productID)
	if err != nil {
		logger.LogError("Error convert product id string to ObjectID")
		return false, http.StatusInternalServerError, err
	}
	for i := 0; i < totalItems; i++ {
		item := &entity.ProductItem{
			ProductID: pID,
			BaseModel: entity.BaseModel{
				Status: common.StatusActive,
			},
			ItemIndex: startIndex + i,
		}
		item.SetTime()

		_, err = s.repo.InsertProductItem(item)
		if err != nil {
			logger.LogError("Get error when inserting productItem" + err.Error())
			return false, http.StatusInternalServerError, err
		}
	}
	return true, http.StatusOK, nil
}

func (s *Service) CountNumProductItems(productID *string) (int, int, error) {
	totalItems, err := s.repo.CountNumProductItems(productID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return totalItems, http.StatusOK, nil
}

func (s *Service) GetProductItemsInOrg(orgID *string) (*[]entity.ProductItem, int, error) {
	pItems, err := s.repo.GetProductItemsInOrg(orgID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return pItems, http.StatusOK, nil
}

func (s *Service) GetProductItemProductOrgAggregate(pID *string) (*entity.ProductItemProductOrgAggregate, int, error) {
	aggregation, err := s.repo.GetProductItemProductOrgAggregate(pID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	aggregation.Product.ParseAttribute()
	return aggregation, http.StatusOK, nil
}
