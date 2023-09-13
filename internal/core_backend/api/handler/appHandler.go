package handler

import (
	"errors"
	"net/http"

	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"

	"backend-service/internal/core_backend/common"

	"github.com/gin-gonic/gin"
)

// AppHandler app handler
type AppHandler struct {
	DummyHandler
	DigitalAssetHandler
	UserHandler
	ScanHandler
	MappingHandler
	ProductItemHandler
	TagHandler
	TemplateHandler
	WebPageHandler
	ProductHandler
	OrganizationHandler
	SessionHandler
	UploadHandler
	PubsubHandler
	AuthorHandler
}

func CreateResponse(err error, code int, xRequestID string, errorMessage string, result interface{}) APIResponse {
	if err != nil {
		logger.LogRequestError(xRequestID, err.Error())
	}

	res := HandlerResponse(code, xRequestID, errorMessage, result)

	return res
}

// APIResponse struct api response
type APIResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	XRequestID string      `json:"x-request-id"`
	Result     interface{} `json:"result"`
}

type RedirectResponse struct {
	StatusCode int    `json:"status_code"`
	URL        string `json:"url"`
}

// HandlerResponse handle response by error code
func HandlerResponse(code int, xRequestID string, errorMessage string, result interface{}) APIResponse {
	if len(errorMessage) == 0 {
		errorMessage = getMessageFromCode(code)
	}
	return APIResponse{
		Code:       code,
		Message:    errorMessage,
		XRequestID: xRequestID,
		Result:     result,
	}
}

// getMessageFromCode get message
func getMessageFromCode(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "Invalid Parameters"
	case common.InvalidStatusCode:
		return "Invalid status"
	case common.NotFoundCode:
		return "Record not found"
	default:
		return http.StatusText(code)
	}
}

func GetRoleFromGinContext(c *gin.Context) (string, error) {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return "", errors.New("userInfo (set at middleware) doesn't exist in Gin Context")
	}

	info := decodeToken.(*entity.User)
	return info.Role, nil
}

func GetUserRoleAndOrgTagNameFromGinContext(c *gin.Context) (string, string, error) {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return "", "", errors.New("userInfo (set at middleware) doesn't exist in Gin Context")
	}
	info := decodeToken.(*entity.User)
	return info.Role, info.Organization, nil
}
