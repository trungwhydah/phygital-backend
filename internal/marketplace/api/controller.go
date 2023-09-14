package api

import "github.com/gin-gonic/gin"

// Controller define interface controller for all controller of http api.
type Controller interface {
	RegisterRoutes(route gin.IRoutes)
}

// RegisterRoutes register the routes of controller class.
func RegisterRoutes(router gin.IRoutes, controllers ...Controller) {
	for _, item := range controllers {
		item.RegisterRoutes(router)
	}
}
