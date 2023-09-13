package handler

import (
	"net/http"

	validation "backend-service/internal/core_backend/infrastructure/validator"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/session"
)

// SessionHandler interface
type SessionHandler interface {
	VerifySession(*gin.Context) APIResponse
}

// sessionHandler struct
type sessionHandler struct {
	SessionService   session.UseCase
	SessionPresenter presenter.ConvertSession
	Validator        validation.CustomValidator
}

// NewSessionHandler create handler
func NewSessionHandler(ss session.UseCase, sp presenter.ConvertSession, v validation.CustomValidator) SessionHandler {
	return &sessionHandler{
		SessionService:   ss,
		SessionPresenter: sp,
		Validator:        v,
	}
}

func (h *sessionHandler) VerifySession(c *gin.Context) APIResponse {
	var request request.SessionInteractRequest
	request.SessionID = c.Query("sessionId")

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	message, code, err := h.SessionService.CheckValidSession(&request.SessionID)
	if err != nil {
		return CreateResponse(err, code, "", "", nil)
	}

	result := h.SessionPresenter.ResponseVerifySession(&request.SessionID, &message)

	return APIResponse{Code: code, Result: result}
}
