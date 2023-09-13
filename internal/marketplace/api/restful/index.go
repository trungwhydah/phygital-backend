package restful

import (
	controller "backend-service/internal/marketplace/api"
	"backend-service/internal/marketplace/api/restful/v1/admin"
	"backend-service/internal/marketplace/api/restful/v1/public"
	"backend-service/internal/marketplace/api/restful/v1/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func registerRoutes(router gin.IRoutes, controllers ...controller.Controller) {
	for _, item := range controllers {
		item.RegisterRoutes(router)
	}
}

var Module = fx.Options(
	// Controllers
	admin.Module,
	user.Module,
	public.Module,

	// Invoke all controllers to register to user router group
	fx.Invoke(fx.Annotate(registerRoutes, fx.ParamTags(`name:"admin-router"`, `group:"admin-controller"`))),
	fx.Invoke(fx.Annotate(registerRoutes, fx.ParamTags(`name:"user-router"`, `group:"user-controller"`))),
	fx.Invoke(fx.Annotate(registerRoutes, fx.ParamTags(`name:"public-router"`, `group:"public-controller"`))),
)
