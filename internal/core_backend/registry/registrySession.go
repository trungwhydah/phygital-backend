package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/session"
)

// Session API
// NewSessionRepository new session repository
func (i *interactor) NewSessionRepository() *repository.SessionRepository {
	return repository.NewSessionRepository(i.mongo)
}

// NewSessionService new session service
func (i *interactor) NewSessionService() *session.Service {
	return session.NewService(i.NewSessionRepository())
}

// NewSessionPresenter
func (i *interactor) NewSessionPresenter() presenter.ConvertSession {
	return presenter.NewPresenterSession()
}

// NewSessionHandler
func (i *interactor) NewSessionHandler() handler.SessionHandler {
	return handler.NewSessionHandler(i.NewSessionService(), i.NewSessionPresenter(), i.NewCustomValidator())
}
