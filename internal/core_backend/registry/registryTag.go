package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/tag"
)

// Tag API
// NewTagRepository new tag repository
func (i *interactor) NewTagRepository() *repository.TagRepository {
	return repository.NewTagRepository(i.mongo)
}

// NewTagService new tag service
func (i *interactor) NewTagService() *tag.Service {
	return tag.NewService(i.NewTagRepository())
}

// NewTagHandler
func (i *interactor) NewTagHandler() handler.TagHandler {
	return handler.NewTagHandler(i.NewTagService(), i.NewMappingService(), i.NewCustomValidator())
}
