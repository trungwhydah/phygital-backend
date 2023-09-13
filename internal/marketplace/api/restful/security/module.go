package security

import (
	"backend-service/internal/marketplace/api/restful/security/authen"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(authen.NewAuthenticatorDecoder),
	fx.Provide(fx.Annotate(authen.NewUserAuthenticator, fx.ResultTags(`name:"user"`))),
	fx.Provide(fx.Annotate(authen.NewAdminAuthenticator, fx.ResultTags(`name:"admin"`))),
)
