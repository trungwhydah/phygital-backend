package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/author"
)

// NewAuthorRepository - new author repository
func (i *interactor) NewAuthorRepository() *repository.AuthorRepository {
	return repository.NewAuthorRepository(i.mongo)
}

// NewAuthorService - new author service
func (i *interactor) NewAuthorService() *author.Service {
	return author.NewService(i.NewAuthorRepository())
}

// NewPresenterAuthor - new presenter author
func (i *interactor) NewPresenterAuthor() presenter.ConvertAuthor {
	return presenter.NewPresenterAuthor()
}

// NewAuthorHandler - New Author Handler
func (i *interactor) NewAuthorHandler() handler.AuthorHandler {
	return handler.NewAuthorHandler(i.NewAuthorService(), i.NewProductService(), i.NewPresenterAuthor())
}
