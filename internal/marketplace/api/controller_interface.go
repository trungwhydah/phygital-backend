package api

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoutes(route gin.IRoutes)
}
