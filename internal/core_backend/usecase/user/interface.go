package user

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User interface
type User interface {
	// Interface for repository
	CreateUser(*entity.User) (*entity.User, error)
	CheckExistedEmail(email *string) (bool, error)
	GetUserByEmail(email *string) (*entity.User, error)
	GetUserByID(userID *string) (*entity.User, error)
	UpsertUser(*entity.User) error
	UpdateRole(role *string, userID *string) (bool, error)
	UpdateOrgID(orgID *primitive.ObjectID, userID *string) (bool, error)
	UpdateUserDetails(*string, *request.UpdateUserDetailsRequest) (bool, error)
	GetUserWithNoWallet() (*[]entity.User, error)
}

// Repository interface
type Repository interface {
	User
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	RegisterUser(*request.CreateUserRequest) (*entity.User, int, error)
	TokenToUser(token *string) (*entity.User, int, error)
	GetUserByID(userID *string) (*entity.User, int, error)
	GetUserByEmail(email *string) (*entity.User, int, error)
	UpsertUserFromFireBase(*string) (*entity.User, int, error)
	UpdateRole(*request.UpdateRoleRequest) (bool, int, error)
	UpdateOrgID(*request.UpdateOrgRequest) (bool, int, error)
	UpdateUserDetails(*string, *request.UpdateUserDetailsRequest) (bool, int, error)
	GetUserWithNoWallet() (*[]entity.User, int, error)
}
