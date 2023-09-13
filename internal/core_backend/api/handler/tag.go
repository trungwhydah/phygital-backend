package handler

import (
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/tag"

	"github.com/gin-gonic/gin"
)

// TagHandler interface
type TagHandler interface {
	CreateTag(*gin.Context) APIResponse
}

// tagHandler struct
type tagHandler struct {
	TagService     tag.UseCase
	MappingService mapping.UseCase
	Validator      validation.CustomValidator
}

// NewTagHandler create handler
func NewTagHandler(tuc tag.UseCase, muc mapping.UseCase, v validation.CustomValidator) TagHandler {
	return &tagHandler{
		TagService:     tuc,
		MappingService: muc,
		Validator:      v,
	}
}

// CreateTag	godoc
// CreateTag	API
//
//	@Summary		Create Tag
//	@Description	Create tag
//	@Tags			tag
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/tag/create [post]
//	@Param			create_tag_request	formData	request.CreateTagRequest	true	"Create Tag Request"
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		400					{object}	APIResponse
func (h *tagHandler) CreateTag(c *gin.Context) APIResponse {
	var request request.CreateTagRequest
	// TODO: Get Token from header and add information to request
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	ok, code, err := h.TagService.CreateTag(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if !ok {
		return CreateResponse(err, code, "", common.MessageErrorCreateTagFail, ok)
	}

	res, code, err := h.MappingService.InitMapping(&request.TagID, &request.OrganizationID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), res)
	}

	return HandlerResponse(code, "", "", res)
}
