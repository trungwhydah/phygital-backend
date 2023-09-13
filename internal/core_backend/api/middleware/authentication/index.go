package authentication

import (
	"net/http"

	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"backend-service/internal/core_backend/infrastructure/repository"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

const (
	USER_INFO_KEY = "userInfo"
)

type Authenticator interface {
	Authenticate(c *gin.Context)
}

type AuthenticationService struct {
	AdminAuth  Authenticator
	UserAuth   Authenticator
	PubsubAuth Authenticator
}

func NewAuthenticationService(fbClient *firebase.FirebaseClient, productItemRepo *repository.ProductItemRepository, productRepo *repository.ProductRepository) *AuthenticationService {
	return &AuthenticationService{
		AdminAuth:  NewAdminAuthenticator(fbClient),
		UserAuth:   NewUserAuthenticator(fbClient),
		PubsubAuth: NewPubsubAuthenticator(),
	}
}

type AuthenticatorDecoder struct {
	fbAuth *auth.Client
}

func NewAuthenticatorDecoder(fbAuth *auth.Client) *AuthenticatorDecoder {
	return &AuthenticatorDecoder{fbAuth: fbAuth}
}

func ResponseUnauthorized(c *gin.Context, err string) {
	logger.LogError(err)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err})
	c.Abort()
}
