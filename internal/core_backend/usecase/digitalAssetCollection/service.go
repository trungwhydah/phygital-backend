package digitalAssetCollection

import (
	"net/http"

	"backend-service/internal/core_backend/entity"
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

// GetCollectionByOrgID
func (s *Service) GetCollectionByOrgID(orgID *string) (*entity.DigitalAssetCollection, int, error) {
	digitalAssetCollection, err := s.repo.GetCollectionByOrgID(orgID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return digitalAssetCollection, http.StatusOK, nil
}

func (s *Service) GetCollectionByID(cID *string) (*entity.DigitalAssetCollection, int, error) {
	dac, err := s.repo.GetCollectionByID(cID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return dac, http.StatusOK, nil
}
