package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             string             `json:"id,omitempty" bson:"_id,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email"`
	EmailVerified  bool               `json:"email_verified"`
	Name           string             `json:"full_name,omitempty" bson:"full_name"`
	Picture        string             `json:"picture"`
	Organization   string             `json:"organization"`
	OrganizationID primitive.ObjectID `json:"org_id" bson:"org_id"`
	Role           string             `json:"role,omitempty" bson:"role"`
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Firebase       Firebase           `json:"firebase"`
	WalletAddress  string             `json:"wallet_address" bson:"wallet_address"`
}
type Firebase struct {
	Identities struct {
		GoogleCom []string `json:"google.com"`
		Email     []string `json:"email"`
	} `json:"identities"`
	SignInProvider string `json:"sign_in_provider"`
}

// CollectionName Collection name of User
func (User) CollectionName() string {
	return "users"
}
func (bm *User) SetTime() *User {
	if bm.CreatedAt.IsZero() {
		bm.CreatedAt = time.Now()
	}
	bm.UpdatedAt = time.Now()

	return bm
}

type UserRole string

const (
	SUPER_ADMIN_ROLE UserRole = "SUPER_ADMIN"
	ORG_ADMIN_ROLE   UserRole = "ORG_ADMIN"
)

func (User) NewUserByFsUser(fsUser *FbUser) User {

	fb := Firebase{
		SignInProvider: fsUser.Provider,
	}

	newUser := User{
		ID:           fsUser.GetID(),
		Name:         fsUser.Username,
		Email:        fsUser.Email,
		Role:         fsUser.Role,
		Status:       "Active",
		Organization: fsUser.Organization,
		Picture:      fsUser.PhotoURL,
		Firebase:     fb,
	}
	newUser.SetTime()
	return newUser
}

type FbUser struct {
	ID            string `json:"id" firestore:"id"`
	Email         string `json:"email" firestore:"email"`
	Username      string `json:"username" firestore:"username"`
	Gender        string `json:"gender,omitempty" firestore:"gender,omitempty" query:"false"`
	Role          string `json:"role,omitempty" firestore:"role,omitempty"`
	DeviceToken   string `json:"-" firestore:"deviceToken,omitempty"`
	LoginType     string `json:"loginType,omitempty" firestore:"loginType,omitempty"`
	OS            string `json:"-" firestore:"os,omitempty"`
	EmailVerified bool   `json:"email_verified" firestore:"email_verified,omitempty"`
	Organization  string `json:"organization,omitempty" firestore:"organization,omitempty"`
	PhotoURL      string `json:"photoURL" firestore:"photoURL,omitempty"`
	Provider      string `json:"provider,omitempty" firestore:"provider"`
}

func (e *FbUser) GetID() string {
	return e.ID
}

func (e *FbUser) SetID(id string) {
	e.ID = id
}

func (e *FbUser) GetEmail() string {
	return e.Email
}
