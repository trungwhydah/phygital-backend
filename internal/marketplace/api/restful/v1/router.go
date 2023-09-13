// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	config "backend-service/config/marketplace"
	"backend-service/internal/marketplace/api/restful/security"
	"backend-service/internal/marketplace/api/restful/security/authen"
	"backend-service/internal/marketplace/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type RouteParams struct {
	fx.In
	Engine    *gin.Engine
	Cfg       *config.Config
	UserAuth  authen.AuthenticatorInterface `name:"user"`
	AdminAuth authen.AuthenticatorInterface `name:"admin"`
}

type Router struct {
	fx.Out
	PublicGroup gin.IRoutes `name:"public-router"`
	AdminGroup  gin.IRoutes `name:"admin-router"`
	UserGroup   gin.IRoutes `name:"user-router"`
}

// @title          Matketplace Backend API
// @version        1.0
// @description    This is a sample server.
// @termsOfService http://swagger.io/terms/

// @BasePath /api/v1

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @query.collection.format multi

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                JWT Token

func NewRouter(params RouteParams, config *config.Config) Router {
	engine := params.Engine

	// Swagger
	if params.Cfg.App.Env != "production" {
		swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
		engine.Use(configSwagger).GET("/swagger/*any", swaggerHandler)
	}

	engine.GET(
		"/",
		func(c *gin.Context) {
			c.String(
				http.StatusOK,
				"Server version: %s",
				config.Version,
			)
		},
	)

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// configs cors
	engine.Use(security.CorsConfigs())

	// public group
	publicGroup := engine.Group("/api/v1")

	// admin group
	adminGroup := publicGroup.Group("/admin").Use(params.AdminAuth.Authenticate)

	// user group
	// userGroup := publicGroup.Group("/user").Use(params.UserAuth.Authenticate)
	userGroup := publicGroup.Group("/user")

	return Router{
		PublicGroup: publicGroup,
		AdminGroup:  adminGroup,
		UserGroup:   userGroup,
	}
}

func configSwagger(c *gin.Context) {
	docs.SwaggerInfo.Host = c.Request.Host
}
