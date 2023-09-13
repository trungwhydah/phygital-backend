package registry

import (
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/nft"
)

// NFT API
// NewNFTRepository new dummy repository
func (i *interactor) NewNFTRepository() *repository.NFTRepository {
	return repository.NewNFTRepository(i.mongo)
}

// NewNFTService new dummy service
func (i *interactor) NewNFTService() *nft.Service {
	return nft.NewService(i.NewNFTRepository())
}
