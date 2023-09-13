package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/template"
)

func (i *interactor) NewTemplateRepository() *repository.TemplateRepository {
	return repository.NewTemplateRepository(i.mongo)
}

func (i *interactor) NewTemplateService() *template.Service {
	return template.NewService(i.NewTemplateRepository())
}

func (i *interactor) NewTemplatePresenter() presenter.ConvertTemplate {
	return presenter.NewPresenterTemplate()
}

func (i *interactor) NewTemplateHandler() handler.TemplateHandler {
	return handler.NewTemplateHandler(i.NewTemplateService(), i.NewTemplatePresenter(), i.NewCustomValidator())
}
