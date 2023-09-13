package registry

import (
	"backend-service/internal/core_backend/infrastructure/repository"
	"backend-service/internal/core_backend/usecase/verification"
)

// Verification API
// NewVerificationRepository new verification repository
func (i *interactor) NewVerificationRepository() *repository.VerificationRepository {
	return repository.NewVerificationRepository(i.mongo)
}

// NewVerificationService new verification service
func (i *interactor) NewVerificationService() *verification.Service {
	return verification.NewService(i.NewVerificationRepository())
}
