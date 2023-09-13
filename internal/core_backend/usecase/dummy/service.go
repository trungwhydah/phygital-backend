package dummy

import (
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
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

// GetDummy get dummy info
func (s *Service) GetDummy(dummyRquest *request.DummyRequest) (*entity.Dummy, int, error) {
	dummy, err := s.repo.GetDummyByID(&dummyRquest.ID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return dummy, http.StatusOK, nil
}
