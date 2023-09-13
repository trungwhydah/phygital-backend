package repo

import (
	"backend-service/internal/marketplace/repo/product"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(product.NewMongoRepo),
)
