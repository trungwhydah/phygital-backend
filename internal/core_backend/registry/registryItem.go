package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/productItem"
)

// Item API
// NewItemRepository new productItem repository
func (i *interactor) NewProductItemRepository() *repository.ProductItemRepository {
	return repository.NewProductItemRepository(i.mongo)
}

// NewItemService new productItem service
func (i *interactor) NewProductItemService() *productItem.Service {
	return productItem.NewService(i.NewProductItemRepository())
}

// NewProductItemPresenter
func (i *interactor) NewProductItemPresenter() presenter.ConvertProductItem {
	return presenter.NewPresenterProductItem()
}

// NewItemHandler
func (i *interactor) NewProductItemHandler() handler.ProductItemHandler {
	return handler.NewProductItemHandler(i.NewUserService(), i.NewProductService(), i.NewProductItemService(), i.NewProductItemPresenter(), i.NewMappingService(), i.NewOrganizationService(), i.NewTemplateService(), i.NewWebPageService(), i.NewCustomValidator(), i.NewDigitalAssetService(), i.NewDigitalAssetCollectionService(), i.NewNFTService())
}
