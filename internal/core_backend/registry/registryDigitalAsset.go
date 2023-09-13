package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/digitalAsset"
)

// DigitalAsset API
// NewDigitalAssetRepository new digitalAsset repository
func (i *interactor) NewDigitalAssetRepository() *repository.DigitalAssetRepository {
	return repository.NewDigitalAssetRepository(i.mongo)
}

// NewDigitalAssetService new digitalAsset service
func (i *interactor) NewDigitalAssetService() *digitalAsset.Service {
	return digitalAsset.NewService(i.NewDigitalAssetRepository())
}

// NewDigitalAssetPresenter
func (i *interactor) NewDigitalAssetPresenter() presenter.ConvertDigitalAsset {
	return presenter.NewPresenterDigitalAsset()
}

// NewDigitalAssetHandler
func (i *interactor) NewDigitalAssetHandler() handler.DigitalAssetHandler {
	return handler.NewDigitalAssetHandler(i.NewDigitalAssetCollectionService(), i.NewDigitalAssetService(), i.NewMappingService(), i.NewUserService(), i.NewProductItemService(), i.NewProductService(), i.NewOrganizationService(), i.NewTemplateService(), i.NewDigitalAssetPresenter(), i.NewCustomValidator())
}
