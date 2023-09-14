package authentication

import (
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminAuthenticator struct {
	fbClient *firebase.FirebaseClient
}

func NewAdminAuthenticator(fbClient *firebase.FirebaseClient) *AdminAuthenticator {
	return &AdminAuthenticator{
		fbClient: fbClient,
	}
}

func (a *AdminAuthenticator) Authenticate(c *gin.Context) {
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
	if !strings.EqualFold(token.Role, string(entity.SUPER_ADMIN_ROLE)) && !strings.EqualFold(token.Role, string(entity.ORG_ADMIN_ROLE)) {
		ResponseUnauthorized(c, "Invalid user role: required super admin or org admin role")
		return
	}
	c.Set(USER_INFO_KEY, a.fbClient.FromTokenToUser(token))
	c.Next()
}
