package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/user"
)

// User API
// NewUserRepository new user repository
func (i *interactor) NewUserRepository() *repository.UserRepository {
	return repository.NewUserRepository(i.mongo)
}

// NewUserService new user service
func (i *interactor) NewUserService() *user.Service {
	return user.NewService(i.firebase, i.NewUserRepository(), i.NewOrganizationRepository())
}

// NewUserPresenter
func (i *interactor) NewUserPresenter() presenter.ConvertUser {
	return presenter.NewPresenterUser()
}

// NewUserHandler
func (i *interactor) NewUserHandler() handler.UserHandler {
	return handler.NewUserHandler(i.NewUserService(), i.NewOrganizationService(), i.NewUserPresenter(), i.NewCustomValidator())
}
