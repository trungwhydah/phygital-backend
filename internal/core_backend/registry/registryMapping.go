package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/mapping"
)

// Mapping API
// NewMappingRepository new mapping repository
func (i *interactor) NewMappingRepository() *repository.MappingRepository {
	return repository.NewMappingRepository(i.mongo)
}

// NewMappingService new mapping service
func (i *interactor) NewMappingService() *mapping.Service {
	return mapping.NewService(i.NewMappingRepository())
}

// NewMappingPresenter
func (i *interactor) NewMappingPresenter() presenter.ConvertMapping {
	return presenter.NewPresenterMapping()
}

// NewMappingHandler
func (i *interactor) NewMappingHandler() handler.MappingHandler {
	return handler.NewMappingHandler(i.NewUserService(), i.NewOrganizationService(), i.NewTagService(), i.NewProductService(), i.NewProductItemService(), i.NewTemplateService(), i.NewMappingService(), i.NewMappingPresenter(), i.NewCustomValidator(), i.NewDigitalAssetService())
}
