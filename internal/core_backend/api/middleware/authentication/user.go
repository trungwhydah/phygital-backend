package authentication

import (
	"backend-service/internal/core_backend/infrastructure/firebase"

	"github.com/gin-gonic/gin"
)

type UserAuthenticator struct {
	fbClient *firebase.FirebaseClient
}

func NewUserAuthenticator(fbClient *firebase.FirebaseClient) *UserAuthenticator {
	return &UserAuthenticator{
		fbClient: fbClient,
	}
}

func (a *UserAuthenticator) Authenticate(c *gin.Context) {
	tokenString, err := a.fbClient.ExtractToken(c)
	if err != nil {
		ResponseUnauthorized(c, err.Error())
		return
	}
	token, err := a.fbClient.VerifyToken(tokenString)
	if err != nil {
		ResponseUnauthorized(c, "Cannot decode token data")
		return
	}
	c.Set(USER_INFO_KEY, a.fbClient.FromTokenToUser(token))
	c.Next()
}
