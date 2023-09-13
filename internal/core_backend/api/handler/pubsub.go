package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/pubsub/v1"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/infrastructure/callers"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/user"

	"github.com/gin-gonic/gin"
)

// PubsubHandler interface
type PubsubHandler interface {
	UpsertUser(*gin.Context) APIResponse
}

// pubsubHandler struct
type pubsubHandler struct {
	UserService user.UseCase
	Validator   validation.CustomValidator
}

// NewPubsubHandler create handler
func NewPubsubHandler(uc user.UseCase, v validation.CustomValidator) PubsubHandler {
	return &pubsubHandler{
		UserService: uc,
		Validator:   v,
	}
}

// UpsertUser	godoc
// UpsertUser	API
//
//	@Summary		Upsert User
//	@Description	Upsert user
//	@Tags			pubsub
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/pubsub/upsert-user [post]
//	@Param			pubsub_received_message	body		interface{}	true	"Pubsub Received Message"
//	@Success		200						{object}	APIResponse{result=entity.User}
//	@Failure		400						{object}	APIResponse
func (h *pubsubHandler) UpsertUser(c *gin.Context) APIResponse {
	var msg pubsub.ReceivedMessage

	if err := c.ShouldBindJSON(&msg); err != nil {
		return CreateResponse(errors.New("Can not receive message"), http.StatusBadRequest, "", "", nil)
	}
	if msg.Message == nil {
		return CreateResponse(errors.New("Can not get message"), http.StatusBadRequest, "", "", nil)
	}

	userID := msg.Message.Attributes["userId"]
	if userID == "" || &userID == nil {
		return CreateResponse(errors.New("userId is required"), http.StatusBadRequest, "", "", nil)
	}
	user, code, err := h.UserService.GetUserByID(&userID)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	if user == nil {
		user, code, err = h.UserService.UpsertUserFromFireBase(&userID)
		if err != nil {
			return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
		}
	}
	// responseData, err := h.requestUserWallet(&user.Email)
	// if err != nil {
	// 	return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	// }
	//
	// updateReq := request.UpdateUserDetailsRequest{
	// 	WalletAddress: &responseData.Result.WalletAddress,
	// }
	// ok, code, err := h.UserService.UpdateUserDetails(&user.ID, &updateReq)
	// if err != nil {
	// 	return CreateResponse(err, code, "", err.Error(), nil)
	// }
	// if !ok {
	// 	err = errors.New("Couldn't update user wallet address")
	// 	return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	// }

	return APIResponse{
		Code:   code,
		Result: user,
	}
}

func (h *pubsubHandler) requestUserWallet(email *string) (*presenter.CreateWalletResponse, error) {
	data := map[string]interface{}{
		"client_id": config.C.Wallet.WALLET_CLIENT_ID,
		"email":     *email,
	}

	req, err := callers.CreateRequest(config.C.Wallet.WALLET_DOMAIN_V1+"/wallets", "POST", data)
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
	responseData := presenter.CreateWalletResponse{}

	// Decode JSON response into responseData
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}
