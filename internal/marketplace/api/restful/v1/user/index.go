package user

import (
	"backend-service/internal/marketplace/api/restful/v1/user/product"
	depinjection "backend-service/pkg/common/dep_injection"
)

var Module = depinjection.BulkProvide(
	[]any{
		product.NewController,
	},
	"user-controller",
)
