package author

import (
	"errors"
	"net/http"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/common/logger"
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

// CreateAuthor - Create author service
func (s *Service) CreateAuthor(author *entity.Author) (*entity.Author, int, error) {
	author.SetTime().SetStatus(common.StatusActive)
	insertedAuthor, err := s.repo.CreateAuthor(author)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return insertedAuthor, http.StatusOK, nil
}

// GetListAuthor -
func (s *Service) GetListAuthor() (*[]entity.Author, int, error) {
	listAuthor, err := s.repo.GetListAuthor()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return listAuthor, http.StatusOK, nil
}

// GetAuthorDetail -
func (s *Service) GetAuthorDetail(authorID *string) (*entity.Author, int, error) {
	author, err := s.repo.GetAuthorByID(authorID)
	if err != nil {
		logger.LogError("")
		return nil, http.StatusInternalServerError, err
	}

	if author == nil {
		return nil, http.StatusOK, errors.New("author not found")
	}

	return author, http.StatusOK, nil
}
