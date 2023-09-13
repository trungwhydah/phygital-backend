package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/scan"
)

// NewScanRepository new scan repository
func (i *interactor) NewScanRepository() *repository.ScanRepository {
	return repository.NewScanRepository(i.mongo)
}

// NewScanService new scan service
func (i *interactor) NewScanService() *scan.Service {
	return scan.NewService(i.caller, i.NewScanRepository())
}

// NewScanPresenter
func (i *interactor) NewScanPresenter() presenter.ConvertScan {
	return presenter.NewPresenterScan()
}

// NewScanHandler
func (i *interactor) NewScanHandler() handler.ScanHandler {
	return handler.NewScanHandler(i.NewSessionService(), i.NewOrganizationService(), i.NewVerificationService(), i.NewProductService(), i.NewProductItemService(), i.NewTagService(), i.NewScanService(), i.NewMappingService(), i.NewTemplateService(), i.NewScanPresenter(), i.NewCustomValidator())
}
