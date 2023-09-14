package handler

import (
	"backend-service/internal/marketplace/api"
	"backend-service/internal/marketplace/api/handler/admin"
	"backend-service/internal/marketplace/api/handler/public"
	"backend-service/internal/marketplace/api/handler/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Controllers
	admin.Module,
	user.Module,
	public.Module,

	// Invoke all controllers to register to user router group
	fx.Invoke(fx.Annotate(api.RegisterRoutes, fx.ParamTags(`name:"admin-router"`, `group:"admin-controller"`))),
	fx.Invoke(fx.Annotate(api.RegisterRoutes, fx.ParamTags(`name:"user-router"`, `group:"user-controller"`))),
	fx.Invoke(fx.Annotate(api.RegisterRoutes, fx.ParamTags(`name:"public-router"`, `group:"public-controller"`))),
)
