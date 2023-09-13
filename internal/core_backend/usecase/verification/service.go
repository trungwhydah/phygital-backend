package verification

import (
	"net/http"
	"strings"

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

// Verify
func (s *Service) Verify(scan *entity.Scan, tag *entity.Tag) (*entity.Verification, int, error) {
	var ver = entity.Verification{
		IsValid: true,
		TagID:   tag.TagID,
	}

	if scan != nil {
		ver.Nonce = scan.ScanCounter
		splitTagID := strings.Split(tag.TagID, "-")
		// 0004 is prefix tag for Ortho
		if splitTagID[0] == "0004" && tag.ScanCounter >= scan.ScanCounter {
			ver.IsValid = false
		}
	} else {
		ver.Nonce = tag.ScanCounter + 1
	}

	ver.SetTime()
	insertedVer, err := s.repo.SaveVerifition(&ver)
	if err != nil {
		logger.LogError("[Debug] - Got error while saving verification" + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return insertedVer, http.StatusOK, nil
}
