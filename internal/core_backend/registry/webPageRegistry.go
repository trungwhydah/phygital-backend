package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/webPage"
)

// NewChipRepository new chip repository
func (i *interactor) NewWebPageRepository() *repository.WebPageRepository {
	return repository.NewWebPageRepository(i.mongo)
}

// NewChipService new chip service
func (i *interactor) NewWebPageService() *webpage.Service {
	return webpage.NewService(i.NewWebPageRepository())
}

func (i *interactor) NewWebPagePresenter() presenter.ConvertWebpage {
	return presenter.NewPresenterWebpage()
}

// NewChipHandler
func (i *interactor) NewWebPageHandler() handler.WebPageHandler {
	return handler.NewWebPageHandler(i.NewWebPageService(), i.NewCustomValidator(), i.NewWebPagePresenter())
}
