package middleware

import (
	"backend-service/internal/core_backend/api/middleware/authentication"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"backend-service/internal/core_backend/infrastructure/repository"
)

type MidddlewareServices struct {
	AuthenMiddleware *authentication.AuthenticationService
}

func NewMiddlewareServices(fbClient *firebase.FirebaseClient, productItemRepo *repository.ProductItemRepository, productRepo *repository.ProductRepository) MidddlewareServices {
	return MidddlewareServices{
		AuthenMiddleware: authentication.NewAuthenticationService(fbClient, productItemRepo, productRepo),
	}
}
