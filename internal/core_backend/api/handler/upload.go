package handler

import (
	"net/http"

	validation "backend-service/internal/core_backend/infrastructure/validator"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/upload"
)

// UploadHandler interface
type UploadHandler interface {
	Upload(*gin.Context) APIResponse
}

// uploadHandler struct
type uploadHandler struct {
	UploadService   upload.UseCase
	UploadPresenter presenter.ConvertUpload
	Validator       validation.CustomValidator
}

// NewUploadHandler create handler
func NewUploadHandler(ds upload.UseCase, up presenter.ConvertUpload, v validation.CustomValidator) UploadHandler {
	return &uploadHandler{
		UploadService:   ds,
		UploadPresenter: up,
		Validator:       v,
	}
}

// Upload	godoc
// Upload	API
//
//	@Summary		Upload Video/Image
//	@Description	Upload video/image
//	@Tags			upload
//	@Accept			multipart/form-data
//	@Produce		json
//	@Router			/upload [post]
//	@Param			files	formData	file	true	"Files"
//	@Success		200		{object}	APIResponse{result=presenter.UploadResponse}
//	@Failure		400		{object}	APIResponse
func (h *uploadHandler) Upload(c *gin.Context) APIResponse {
	var request request.UploadImagesRequest
	err := c.Request.ParseMultipartForm(32 << 20) // max memory of 32MB
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	request.Images = c.Request.MultipartForm.File["files"]
	resultUpload, code, err := h.UploadService.UploadImagesToGCPStorage(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.UploadPresenter.ResponseUpload(resultUpload)

	return APIResponse{Code: code, Result: &result}
}
