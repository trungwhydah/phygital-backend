package authentication

import (
	"github.com/gin-gonic/gin"

	config "backend-service/config/core_backend"
)

type PubsubAuthenticator struct {
}

func NewPubsubAuthenticator() *PubsubAuthenticator {
	return &PubsubAuthenticator{}
}

func (a *PubsubAuthenticator) Authenticate(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		ResponseUnauthorized(c, "Can not get token")
		return
	}
	//TODO: decode token by key

	if token != config.C.Pubsub.AuthenToken {
		ResponseUnauthorized(c, "Invalid Token")
		return
	}
	c.Next()
}
