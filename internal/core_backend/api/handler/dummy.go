package handler

import (
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/dummy"
)

// DummyHandler interface
type DummyHandler interface {
	GetDummy(*gin.Context) APIResponse
}

// dummyHandler struct
type dummyHandler struct {
	DummyService   dummy.UseCase
	DummyPresenter presenter.ConvertDummy
	Validator      validation.CustomValidator
}

// NewDummyHandler create handler
func NewDummyHandler(ds dummy.UseCase, dp presenter.ConvertDummy, v validation.CustomValidator) DummyHandler {
	return &dummyHandler{
		DummyService:   ds,
		DummyPresenter: dp,
		Validator:      v,
	}
}

func (h *dummyHandler) GetDummy(c *gin.Context) APIResponse {
	var request request.DummyRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	dummy, code, err := h.DummyService.GetDummy(&request)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	result := h.DummyPresenter.ResponseDummy(dummy)

	return APIResponse{Code: code, Result: result}
}
