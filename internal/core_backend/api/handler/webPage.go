package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	webpage "backend-service/internal/core_backend/usecase/webPage"
)

// WebPageHandler interface
type WebPageHandler interface {
	CreateWebPage(*gin.Context) APIResponse
	GetAllWebPages(*gin.Context) APIResponse
	GetWebPage(*gin.Context) APIResponse
	UpdateWebPage(*gin.Context) APIResponse
	DeleteWebPage(*gin.Context) APIResponse
}

// webPageHandler struct
type webPageHandler struct {
	WebPageService   webpage.UseCase
	Validator        validation.CustomValidator
	WebpagePresenter presenter.ConvertWebpage
}

// NewChipHandler create handler
func NewWebPageHandler(wu webpage.UseCase, v validation.CustomValidator, wp presenter.ConvertWebpage) WebPageHandler {
	return &webPageHandler{
		WebPageService:   wu,
		Validator:        v,
		WebpagePresenter: wp,
	}
}

// CreateWebPage	godoc
// CreateWebPage	API
//
//	@Summary		Create Webpage
//	@Description	Create webpage
//	@Tags			webpage
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/web-page/create [post]
//	@Param			create_webpage_request	body		request.CreateWebpageRequest	true	"Create Webpage Request"
//	@Success		200						{object}	APIResponse{result=string}
//	@Failure		500						{object}	APIResponse
func (h *webPageHandler) CreateWebPage(c *gin.Context) APIResponse {
	var request request.CreateWebpageRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}
	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}
	id, code, err := h.WebPageService.CreateWebPage(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", id)
}

// GetAllWebPages	godoc
// GetAllWebPages	API
//
//	@Summary		Get List Of All Webpages
//	@Description	Get list of all webpages
//	@Tags			webpage
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/web-page/all [get]
//	@Success		200	{object}	APIResponse{result=presenter.AllWebpagesResponse}
//	@Failure		500	{object}	APIResponse
func (h *webPageHandler) GetAllWebPages(c *gin.Context) APIResponse {

	pages, code, err := h.WebPageService.GetAllWebPages()
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.WebpagePresenter.ResponseAllWebpages(pages)

	return HandlerResponse(code, "", "", result)
}

// GetWebPages	godoc
// GetWebPages	API
//
//	@Summary		Get Webpage By ID
//	@Description	Get webPage by ID from DB
//	@Tags			webpage user
//	@Produce		json
//	@Router			/web-page/{webpage_id} [get]
//	@Param			webpage_id	path		string	true	"Webpage ID"
//	@Success		200			{object}	APIResponse{result=presenter.WebpageDetailResponse}
//	@Failure		500			{object}	APIResponse
func (h *webPageHandler) GetWebPage(c *gin.Context) APIResponse {

	webpageID := c.Param("webpage_id")
	if webpageID == "" {
		return CreateResponse(errors.New("Invalid webpageID"), http.StatusBadRequest, "", "", nil)
	}

	page, code, err := h.WebPageService.GetWebPage(&webpageID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.WebpagePresenter.ResponseWebpageDetail(page)

	return HandlerResponse(code, "", "", result)
}

// UpdateWebPage	godoc
// UpdateWebPage	API
//
//	@Summary		Update Webpage by ID
//	@Description	Update webpages by id from DB
//	@Tags			webpage
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/web-page/{webpage_id} [put]
//	@Param			webpage_id				path		string							true	"Webpage ID"
//	@Param			update_webpage_request	body		request.UpdateWebpageRequest	true	"Update Webpage Request"
//	@Success		200						{object}	APIResponse{result=string}
//	@Failure		400						{object}	APIResponse
func (h *webPageHandler) UpdateWebPage(c *gin.Context) APIResponse {

	var request request.UpdateWebpageRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}
	request.WebpageID = c.Param("webpage_id")

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}
	id, code, err := h.WebPageService.UpdateWebPage(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", id)
}

// DeleteWebPage	godoc
// DeleteWebPage	API
//
//	@Summary		Delete Webpage
//	@Description	Delete webpage
//	@Tags			webpage
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/web-page/{webpage_id} [delete]
//	@Param			webpage_id	path		string	true	"Webpage ID"
//	@Success		200			{object}	APIResponse{result=bool}
//	@Failure		500			{object}	APIResponse
func (h *webPageHandler) DeleteWebPage(c *gin.Context) APIResponse {
	pageID := c.Param("webpage_id")
	if pageID == "" {
		return CreateResponse(errors.New("Invalid pageID"), http.StatusBadRequest, "", "", nil)
	}

	result, code, err := h.WebPageService.DeleteWebPage(&pageID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), result)
	}

	return HandlerResponse(code, "", "", result)
}
