package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/infrastructure/callers"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/user"

	"github.com/gin-gonic/gin"
)

// UserHandler interface
type UserHandler interface {
	UserSignUp(*gin.Context) APIResponse
	UpdateRole(*gin.Context) APIResponse
	UpdateOrg(*gin.Context) APIResponse
	UpdateUserDetails(*gin.Context) APIResponse
	GetUserDetails(*gin.Context) APIResponse
	SyncWalletAddress(*gin.Context) APIResponse
}

// userHandler struct
type userHandler struct {
	UserService         user.UseCase
	OrganizationService organization.UseCase
	UserPresenter       presenter.ConvertUser
	Validator           validation.CustomValidator
}

// NewUserHandler create handler
func NewUserHandler(uc user.UseCase, ou organization.UseCase, pr presenter.ConvertUser, v validation.CustomValidator) UserHandler {
	return &userHandler{
		UserService:         uc,
		OrganizationService: ou,
		UserPresenter:       pr,
		Validator:           v,
	}
}

func (h *userHandler) UserSignUp(c *gin.Context) APIResponse {
	var request request.CreateUserRequest
	request.Token = c.GetHeader("Authorization")

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	user, code, err := h.UserService.RegisterUser(&request)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	return APIResponse{
		Code:   code,
		Result: h.UserPresenter.ResponseUser(user),
	}
}

// UpdateRole	godoc
// UpdateRole	API
//
//	@Summary		Update User Role To Firebase and DB
//	@Description	Update User Role To Firebase and DB
//	@Tags			user
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/user/role [put]
//	@Param			update_role_request	formData	request.UpdateRoleRequest	true	"Update Role Request"
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		400					{object}	APIResponse
func (h *userHandler) UpdateRole(c *gin.Context) APIResponse {
	userRole, err := GetRoleFromGinContext(c)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}
	if userRole != string(entity.SUPER_ADMIN_ROLE) {
		err = errors.New("Unauthorized: only super admin can update user role")
		return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
	}
	var req request.UpdateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	success, code, err := h.UserService.UpdateRole(&req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	return HandlerResponse(code, "", "", success)
}

// UpdateOrg	godoc
// UpdateOrg	API
//
//	@Summary		Update User Org To Firebase and DB
//	@Description	Update User Org To Firebase and DB
//	@Tags			user
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/user/org [put]
//	@Param			update_org_request	formData	request.UpdateOrgRequest	true	"Update Org Request"
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		400					{object}	APIResponse
func (h *userHandler) UpdateOrg(c *gin.Context) APIResponse {
	var req request.UpdateOrgRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	org, code, err := h.OrganizationService.GetOrgByTagName(&req.NewOrgTagName)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if org == nil {
		err = errors.New("Given organization tag name doesn't exist")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	req.NewOrgID = org.ID

	success, code, err := h.UserService.UpdateOrgID(&req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	return HandlerResponse(code, "", "", success)
}

// UpdateUserDetails	godoc
// UpdateUserDetails	API
//
//	@Summary		Update User Details
//	@Description	Update User Details
//	@Tags			user
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/user [put]
//	@Param			update_user_details_request	body		request.UpdateUserDetailsRequest	true	"Update User Details Request"
//	@Success		200							{object}	APIResponse{result=bool}
//	@Failure		400							{object}	APIResponse
func (h *userHandler) UpdateUserDetails(c *gin.Context) APIResponse {
	var req request.UpdateUserDetailsRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		err := errors.New("userInfo (set at middleware) doesn't exist in Gin Context")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	info := decodeToken.(*entity.User)

	success, code, err := h.UserService.UpdateUserDetails(&info.ID, &req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	return HandlerResponse(code, "", "", success)
}

// GetUserDetails	godoc
// GetUserDetails	API
//
//	@Summary		Get User Details
//	@Description	Get User Details
//	@Tags			user
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/user [get]
//	@Success		200	{object}	APIResponse{result=entity.User}
//	@Failure		400	{object}	APIResponse
func (h *userHandler) GetUserDetails(c *gin.Context) APIResponse {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		err := errors.New("userInfo (set at middleware) doesn't exist in Gin Context")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	info := decodeToken.(*entity.User)

	user, code, err := h.UserService.GetUserByID(&info.ID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	return HandlerResponse(code, "", "", user)
}

// SyncWalletAddress	godoc
// SyncWalletAddress	API
//
//	@Summary		Sync Wallet Address Of Users Who Haven't Had
//	@Description	Sync Wallet Address Of Users Who Haven't Had
//	@Tags			user
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/user/sync-wallet-address [put]
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		400					{object}	APIResponse
func (h *userHandler) SyncWalletAddress(c *gin.Context) APIResponse {
	userRole, err := GetRoleFromGinContext(c)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}
	if userRole != string(entity.SUPER_ADMIN_ROLE) {
		err = errors.New("Unauthorized: only super admin can sync user wallet addresses")
		return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
	}
	usersNoWallet, code, err := h.UserService.GetUserWithNoWallet()
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if len(*usersNoWallet) == 0 {
		return HandlerResponse(code, "", "", true)
	}
	emailToID := make(map[string]string)
	var emails []string
	for _, user := range *usersNoWallet {
		emails = append(emails, user.Email)
		emailToID[user.Email] = user.ID
	}
	responseData, err := h.requestUserWalletBulk(&emails)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}

	// Update wallet_address for each user
	for _, userWallet := range responseData.Result.ListWallet {
		updateReq := request.UpdateUserDetailsRequest{
			WalletAddress: &userWallet.Wallet.WalletAddress,
		}
		userID := emailToID[userWallet.Email]
		ok, code, err := h.UserService.UpdateUserDetails(&userID, &updateReq)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if !ok {
			err = fmt.Errorf("Couldn't update user info: email=%s, _id=%s", userWallet.Email, userID)
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
	}
	return HandlerResponse(code, "", "", true)
}

func (h *userHandler) requestUserWalletBulk(emails *[]string) (*presenter.CreateWalletsBulkResponse, error) {
	data := map[string]interface{}{
		"client_id": config.C.Wallet.WALLET_CLIENT_ID,
		"emails":    emails,
	}

	req, err := callers.CreateRequest(config.C.Wallet.WALLET_DOMAIN_V1+"/wallets/bulk", "POST", data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-client-key", config.C.Wallet.WALLET_CLIENT_KEY)

	client := http.Client{
		Timeout: time.Duration(config.C.Wallet.WALLET_API_TIMEOUT) * time.Second,
	}

	resp, err := callers.SendRequest(&client, req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Received non-OK response: %s", resp.Status)
		return nil, err
	}
	responseData := presenter.CreateWalletsBulkResponse{}

	// Decode JSON response into responseData
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}
