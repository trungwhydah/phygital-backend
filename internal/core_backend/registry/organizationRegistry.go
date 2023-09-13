package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/organization"
)

// Organization API
// NewOrganizationRepository new organization repository
func (i *interactor) NewOrganizationRepository() *repository.OrganizationRepository {
	return repository.NewOrganizationRepository(i.mongo)
}

// NewOrganizationService new organization service
func (i *interactor) NewOrganizationService() *organization.Service {
	return organization.NewService(i.NewOrganizationRepository())
}

// NewOrganizationPresenter
func (i *interactor) NewOrganizationPresenter() presenter.ConvertOrganization {
	return presenter.NewPresenterOrganization()
}

// NewOrganizationHandler
func (i *interactor) NewOrganizationHandler() handler.OrganizationHandler {
	return handler.NewOrganizationHandler(i.NewOrganizationService(), i.NewOrganizationPresenter(), i.NewCustomValidator())
}
