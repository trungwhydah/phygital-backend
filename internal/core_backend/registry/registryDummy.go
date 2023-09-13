package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/dummy"
)

// Dummy API
// NewDummyRepository new dummy repository
func (i *interactor) NewDummyRepository() *repository.DummyRepository {
	return repository.NewDummyRepository(i.mongo)
}

// NewDummyService new dummy service
func (i *interactor) NewDummyService() *dummy.Service {
	return dummy.NewService(i.NewDummyRepository())
}

// NewDummyPresenter
func (i *interactor) NewDummyPresenter() presenter.ConvertDummy {
	return presenter.NewPresenterDummy()
}

// NewDummyHandler
func (i *interactor) NewDummyHandler() handler.DummyHandler {
	return handler.NewDummyHandler(i.NewDummyService(), i.NewDummyPresenter(), i.NewCustomValidator())
}
