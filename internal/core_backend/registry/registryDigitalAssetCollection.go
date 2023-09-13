package registry

import (
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/digitalAssetCollection"
)

// DigitalAssetCollection API
// NewDigitalAssetCollectionRepository new digitalAssetCollection repository
func (i *interactor) NewDigitalAssetCollectionRepository() *repository.DigitalAssetCollectionRepository {
	return repository.NewDigitalAssetCollectionRepository(i.mongo)
}

// NewDigitalAssetCollectionService new digitalAssetCollection service
func (i *interactor) NewDigitalAssetCollectionService() *digitalAssetCollection.Service {
	return digitalAssetCollection.NewService(i.NewDigitalAssetCollectionRepository())
}
