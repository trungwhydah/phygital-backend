package middleware

import (
	"backend-service/internal/marketplace/api/middleware/authen"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(authen.NewAuthenticatorDecoder),
	fx.Provide(fx.Annotate(authen.NewUserAuthenticator, fx.ResultTags(`name:"user"`))),
	fx.Provide(fx.Annotate(authen.NewAdminAuthenticator, fx.ResultTags(`name:"admin"`))),
)
