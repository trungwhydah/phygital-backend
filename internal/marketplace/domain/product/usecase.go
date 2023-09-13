package product

import (
	"context"

	"backend-service/internal/marketplace/entity/product"
	productrepo "backend-service/internal/marketplace/repo/product"
)

type UseCaseInterface interface {
	GetProduct(ctx context.Context, productID *string) (*product.Product, error)
}

type UseCase struct {
	productRepo productrepo.RepoInterface
}

func (u *UseCase) GetProduct(ctx context.Context, productID *string) (*product.Product, error) {
	return u.productRepo.FindOneByID(ctx, productID)
}

func NewUseCase(
	productRepo productrepo.RepoInterface,
) UseCaseInterface {
	return &UseCase{
		productRepo: productRepo,
	}
}
