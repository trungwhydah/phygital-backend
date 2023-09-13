package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateUserRequest struct {
	Token string `form:"token" validate:"required"`
}

type UpdateRoleRequest struct {
	UserID  string `form:"user_id" validate:"required"`
	NewRole string `form:"new_role" validate:"required"`
}

type UpdateOrgRequest struct {
	UserID        string             `form:"user_id" validate:"required"`
	NewOrgTagName string             `form:"new_org_tag_name" validate:"required"`
	NewOrgID      primitive.ObjectID `swaggerignore:"true"`
}

type UpdateUserDetailsRequest struct {
	Name          *string `json:"full_name" bson:"full_name,omitempty"`
	Picture       *string `json:"picture" bson:"picture,omitempty"`
	WalletAddress *string `json:"wallet_address" bson:"wallet_address"`
}
