package tag

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"
	"errors"
	"net/http"

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

// CreateTag
func (s *Service) CreateTag(request *request.CreateTagRequest) (bool, int, error) {
	oID, err := primitive.ObjectIDFromHex(request.OrganizationID)
	if err != nil {
		logger.LogError("gor error when parsing organization")
		return false, http.StatusInternalServerError, err
	}

	tag := &entity.Tag{
		HardwareID:     request.HardwareID,
		TagID:          request.TagID,
		OrganizationID: oID,
		EncryptMode:    request.EncryptMode,
		RawData:        request.RawData,
	}

	isExisted, err := s.repo.CheckExistedTag(tag)
	if err != nil {
		logger.LogError("Get error when checking existed: " + err.Error())
		return false, http.StatusInternalServerError, err
	}
	if isExisted {
		return false, http.StatusOK, errors.New(common.MessageErrorExistedTag)
	}

	tag.SetTime()
	tagInserted, err := s.repo.CreateTag(tag)
	if err != nil {
		logger.LogError("Get error when creating tag: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return !tagInserted.ID.IsZero(), http.StatusOK, nil
}

func (s *Service) GetTagNotMapped(tagMapped *[]string) (*[]entity.Tag, int, error) {
	tags, err := s.repo.GetTagNotMapped(tagMapped)
	if err != nil {
		logger.LogError("Failed to get tags not mapped: " + err.Error())
		return nil, http.StatusInternalServerError, nil
	}

	return tags, http.StatusOK, nil
}

func (s *Service) GetTag(tagID string) (*entity.Tag, int, error) {
	tag, _ := s.repo.GetTag(tagID)
	return tag, http.StatusOK, nil
}

// UpdateTagCounter
func (s *Service) UpdateTagCounter(tagID *string, scanCounter *int) (bool, int, error) {
	success, err := s.repo.UpdateTagCounter(tagID, scanCounter)
	if err != nil {
		logger.LogError("[DEBUG] - 7 - Error updating tag counter: " + err.Error())
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

// GetTagByHWID - Get tag infor by Hardware ID
func (s *Service) GetTagByHWID(uid *string) (*entity.Tag, int, error) {
	tag, err := s.repo.GetTagByHWID(uid)
	if err != nil {
		logger.LogError("[DEBUG] - 6 - Get error while getting tag information: " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return tag, http.StatusOK, nil
}
