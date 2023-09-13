package mapping

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

// GetAllMappingInOrg
func (s *Service) GetAllMappingInOrg(orgID *string) (*[]entity.Mapping, int, error) {
	var mappingList *[]entity.Mapping
	if len(*orgID) != 0 {
		oID, err := primitive.ObjectIDFromHex(*orgID)
		if err != nil {
			logger.LogError("Got error while parsing organization: " + err.Error())
			return nil, http.StatusInternalServerError, err
		}
		mappingList, err = s.repo.GetAllMappingInOrg(&oID)
		if err != nil {
			logger.LogError("Got error while getting mapping: " + err.Error())
			return nil, http.StatusInternalServerError, err
		}
	} else {
		result, err := s.repo.GetAllMapping()
		if err != nil {
			logger.LogError("Got error while getting mapping: " + err.Error())
			return nil, http.StatusInternalServerError, err
		}

		mappingList = result
	}

	return mappingList, http.StatusOK, nil
}

func (s *Service) GetMappingWithTagID(tagID *string) (*entity.Mapping, int, error) {
	mapping, err := s.repo.GetMappingWithTagID(tagID)
	if err != nil {
		logger.LogError("[DEBUG] - 8 - Got error while getting mapping: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return mapping, http.StatusOK, nil
}

func (s *Service) InitMapping(tagID, orgID *string) (bool, int, error) {
	oID, err := primitive.ObjectIDFromHex(*orgID)
	if err != nil {
		logger.LogError("Error convert org id from string to ObjectID")
		return false, http.StatusInternalServerError, err
	}
	mapping := &entity.Mapping{
		OrganizationID: oID,
		TagID:          *tagID,
	}
	mapping.SetTime()
	success, err := s.repo.UpsertMapping(mapping)
	if err != nil {
		logger.LogError("Get error when creating mapping: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) IsProductItemIDMapped(productItemID *string) (bool, int, error) {
	mapping, err := s.repo.GetMappingWithProductItemID(productItemID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	if mapping == nil {
		return false, http.StatusOK, nil
	}
	return true, http.StatusOK, nil
}

func (s *Service) GetMappingWithProductItemID(productItemID *string) (*entity.Mapping, int, error) {
	mapping, err := s.repo.GetMappingWithProductItemID(productItemID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return mapping, http.StatusOK, nil
}

// GetAllMappingForProduct
func (s *Service) GetAllMappingForProduct(productID *string, orgID *string) (*[]entity.Mapping, int, error) {
	var mappingList *[]entity.Mapping
	mappingList, err := s.repo.GetAllMappingForProduct(productID, orgID)
	if err != nil {
		logger.LogError("Got error while getting mapping: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return mappingList, http.StatusOK, nil
}

func (s *Service) UpdateMapping(tagID *string, req *request.UpdateMappingRequest) (bool, int, error) {
	ok, err := s.repo.UpdateMapping(tagID, req)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	return ok, http.StatusOK, nil
}

func (s *Service) Unmap(tagID *string) (bool, int, error) {
	ok, err := s.repo.Unmap(tagID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	return ok, http.StatusOK, nil
}

// GetMappingByDigitalAsset
func (s *Service) GetMappingByDigitalAsset(digitalAssetID *string) (*entity.Mapping, int, error) {
	mapping, err := s.repo.GetMappingByDigitalAsset(digitalAssetID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return mapping, http.StatusOK, nil
}
