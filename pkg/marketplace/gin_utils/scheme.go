package ginutils

import "github.com/gin-gonic/gin"

func Scheme(g *gin.Context) string {
	scheme := "http"
	if g.Request.TLS != nil {
		scheme = "https"
	}

	return scheme
}
