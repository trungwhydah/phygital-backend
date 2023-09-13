package session

import (
	"backend-service/internal/core_backend/entity"
)

// Session interface
type Session interface {
	// Interface for repository
	CreateSession(*entity.Session) (*entity.Session, error)
	GetSessionWithID(sessionID *string) (*entity.Session, error)
}

// Repository interface
type Repository interface {
	Session
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateSession(tagID *string, timeoutInSecond int) (*entity.Session, int, error)
	CheckValidSession(sessionId *string) (string, int, error)
}
