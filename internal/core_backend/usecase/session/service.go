package session

import (
	"net/http"
	"time"

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

// CreateSession Create
func (s *Service) CreateSession(tagID *string, timeoutInSecond int) (*entity.Session, int, error) {
	currentTime := time.Now()
	session := &entity.Session{
		TagID:     *tagID,
		StartAt:   currentTime,
		ExpiredAt: currentTime.Add(time.Duration(timeoutInSecond) * time.Second),
	}
	session, err := s.repo.CreateSession(session)

	if err != nil {
		logger.LogError("[DEBUG] - 9 - Error creating session" + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return session, http.StatusOK, nil
}

// CheckValidSession
func (s *Service) CheckValidSession(sessionId *string) (string, int, error) {
	session, err := s.repo.GetSessionWithID(sessionId)
	if err != nil {
		logger.LogError("got error while getting session " + err.Error())
		return "fail to get session", http.StatusInternalServerError, err
	}

	if session == nil {
		return "session is not found", http.StatusOK, nil
	}

	if session.ExpiredAt.Before(time.Now()) {
		return "session is expired", http.StatusOK, nil
	}

	return "", http.StatusOK, nil
}
