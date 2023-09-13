package security

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

// CorsConfigs sets up cors middleware.
func CorsConfigs() gin.HandlerFunc {
	return cors.New(
		cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE"},
			ExposedHeaders:   []string{"Content-Length"},
			MaxAge:           86400,
			AllowCredentials: true,
		},
	)
}
