package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/product"
)

// Product API
// NewProductRepository new product repository
func (i *interactor) NewProductRepository() *repository.ProductRepository {
	return repository.NewProductRepository(i.mongo)
}

// NewProductService new product service
func (i *interactor) NewProductService() *product.Service {
	return product.NewService(i.NewProductRepository())
}

// NewProductPresenter
func (i *interactor) NewProductPresenter() presenter.ConvertProduct {
	return presenter.NewPresenterProduct()
}

// NewProductHandler
func (i *interactor) NewProductHandler() handler.ProductHandler {
	return handler.NewProductHandler(i.NewMappingService(), i.NewOrganizationService(), i.NewWebPageService(), i.NewProductService(), i.NewProductItemService(), i.NewTemplateService(), i.NewProductPresenter(), i.NewCustomValidator())
}
