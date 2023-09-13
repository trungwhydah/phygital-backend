package firebase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	App *firebase.App
}

type Token struct {
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Role          string `json:"role"`
	Organization  string `json:"organization"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	AuthTime      int    `json:"auth_time"`
	UserID        string `json:"user_id"`
	Sub           string `json:"sub"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Firebase      struct {
		Identities struct {
			GoogleCom []string `json:"google.com"`
			Email     []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
}

func NewFirebaseClient(projectID string) (*FirebaseClient, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: projectID},
		option.WithCredentialsFile("./service-account-file.json"),
	)
	if err != nil {
		logger.LogError("Error when initializing Firebase app: " + err.Error())
		return nil, err
	}

	return &FirebaseClient{
		App: app,
	}, nil
}

func (fc *FirebaseClient) VerifyToken(token string) (*Token, error) {
	ctx := context.Background()
	authApp, err := fc.App.Auth(ctx)
	if err != nil {
		logger.LogError("Error when initializing Authenticate of Firebase app: " + err.Error())
		return nil, err
	}

	res, err := authApp.VerifyIDToken(ctx, token)
	if err != nil {
		logger.LogError("Error when verifing Token: " + err.Error())
		return nil, err
	}

	encodeToken, err := json.Marshal(res.Claims)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error while decoding claims: %v", err.Error()))
		return nil, err
	}

	var tk Token
	err = json.Unmarshal(encodeToken, &tk)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error while unmarshaling claims: %v", err.Error()))
		return nil, err
	}

	return &tk, nil
}

func (fc *FirebaseClient) ExtractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.LogInfo("No authorization header")
		return "", errors.New("No authorization header")
	}
	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "bearer") {
		logger.LogError("Invalid authorization header")
		return "", errors.New("Invalid authorization header")
	}
	return authParts[1], nil
}

func (fc *FirebaseClient) FromTokenToUser(token *Token) *entity.User {
	return &entity.User{
		ID:            token.UserID,
		Email:         token.Email,
		EmailVerified: token.EmailVerified,
		Name:          token.Name,
		Picture:       token.Picture,
		Organization:  token.Organization,
		Role:          token.Role,
		Firebase:      token.Firebase,
	}
}

func (fc *FirebaseClient) GetFirebaseUser(userID string) (*entity.FbUser, error) {
	ctx := context.Background()
	authApp, err := fc.App.Auth(ctx)
	if err != nil {
		logger.LogError("Error when initializing Authenticate of Firebase app: " + err.Error())
		return nil, err
	}

	res, err := authApp.GetUser(ctx, userID)
	if err != nil {
		logger.LogError("Error when get user info: " + err.Error())
		return nil, err
	}
	userInfo := entity.FbUser{}
	userInfo.ID = res.UserInfo.UID
	userInfo.Username = res.UserInfo.DisplayName
	if res.CustomClaims["role"] != nil {
		userInfo.Role = fmt.Sprint(res.CustomClaims["role"])
	}
	if len(res.ProviderUserInfo) > 0 {
		userInfo.Provider = res.ProviderUserInfo[0].ProviderID
	}
	userInfo.PhotoURL = res.UserInfo.PhotoURL
	userInfo.EmailVerified = res.EmailVerified
	userInfo.Email = res.UserInfo.Email
	if res.CustomClaims["organization"] != nil {
		userInfo.Organization = fmt.Sprint(res.CustomClaims["organization"])
	}
	return &userInfo, nil
}

func (fc *FirebaseClient) UpdateUserDetails(user *entity.User) error {
	if user == nil {
		return errors.New("Error getting updated user info")
	}
	ctx := context.Background()
	authApp, err := fc.App.Auth(ctx)
	if err != nil {
		logger.LogError("Error when initializing Authenticate of Firebase app: " + err.Error())
		return err
	}

	// Update the custom claims
	claims := map[string]interface{}{
		"role":         user.Role,
		"organization": user.Organization,
	}

	err = authApp.SetCustomUserClaims(context.Background(), user.ID, claims)
	if err != nil {
		logger.LogError("Failed to set custom claims: " + err.Error())
		return err
	}

	return nil
}
