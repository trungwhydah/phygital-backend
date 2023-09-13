package domain

import (
	"backend-service/internal/marketplace/domain/product"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(product.NewUseCase),
)
