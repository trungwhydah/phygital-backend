package presenter

import (
	"backend-service/internal/core_backend/entity"
)

// UserResponse data struct
type UserResponse struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// presenterUser struct
type PresenterUser struct{}

// presenterUser interface
type ConvertUser interface {
	ResponseUser(user *entity.User) *UserResponse
}

// NewPresenterUser Constructs presenter
func NewPresenterUser() ConvertUser {
	return &PresenterUser{}
}

// Return property data response
func (pp *PresenterUser) ResponseUser(user *entity.User) *UserResponse {
	response := &UserResponse{
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	}

	return response
}
