package httpresp

import (
	"net/http"

	"backend-service/pkg/common/msgtranslate"
	"backend-service/pkg/common/pagination"
	"backend-service/pkg/common/utils"

	"github.com/gin-gonic/gin"
)

// Response define struct returns for http api.
//
// Data specific all data http request.
// Pagination specific pagination information when http request get list items.
// ErrorKey specific error key if http request have error.
// Message specific detail error if http request have error, message will be translated to language based on header of
// http request.
type Response struct {
	Data       any                    `json:"data,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	ErrorKey   *string                `json:"error_key,omitempty" example:"error.system.internal"`
	Message    *string                `json:"message,omitempty" example:"Internal System Error"`
}

// Error returns error for rest api.
//
// status specific HTTP code err.
// errorKey specific for message of err.
// msgArgs specific for dynamic variable for message of err.
func Error(
	c *gin.Context,
	status int,
	errorKey string,
	msgArgs map[string]any,
) {
	lang := GetLanguageCode(c)

	c.AbortWithStatusJSON(
		status,
		Response{
			ErrorKey: &errorKey,
			Message:  utils.ToPtr(msgtranslate.Translate(errorKey, &lang, msgArgs)),
		},
	)
}

// Success returns result for rest api.
//
// res specific result of rest api.
func Success(c *gin.Context, res *Response) {
	c.JSON(http.StatusOK, res)
}

// SuccessNoContent returns success for rest api without content.
func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// InternalServerError returns error result for rest api with internal system error.
func InternalServerError(g *gin.Context) {
	Error(
		g,
		http.StatusInternalServerError,
		ErrKeySystemInternalServer.Error(),
		nil,
	)
}

// UnauthorizedError returns error result for rest api with authentication no permission.
//
// errorMsg specific reason for error don't have permission.
func UnauthorizedError(g *gin.Context, errorMsg string) {
	Error(
		g,
		http.StatusUnauthorized,
		ErrKeyAuthenticationNoPermission.Error(),
		map[string]any{
			"msg_err": errorMsg,
		},
	)
}

// MissingRequiredFieldError returns error result for rest api with.
//
// field specific field missing.
func MissingRequiredFieldError(g *gin.Context, field string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrKeyHTTPValidatorsMissingRequiredField.Error(),
		map[string]any{"field": field},
	)
}

// InvalidFieldTypeError returns error result for rest api with invalid type.
//
// msgErr specific reason for decode fail.
func InvalidFieldTypeError(g *gin.Context, msgErr string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrKeyHTTPValidatorsInvalidFieldType.Error(),
		map[string]any{
			"msg_err": msgErr,
		},
	)
}

// DecodeFail returns error result for rest api with authentication no permission.
//
// msgErr specific reason for decode fail.
func DecodeFail(g *gin.Context, msgErr string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrKeyHTTPValidatorsDecodeFail.Error(),
		map[string]any{
			"msg_err": msgErr,
		},
	)
}

func NotFound(g *gin.Context) {
	Error(
		g,
		http.StatusNotFound,
		ErrKeyDatabaseNotFound.Error(),
		nil,
	)
}
