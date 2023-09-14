package handler

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/common"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TemplateHandler interface {
	CreateTemplate(*gin.Context) APIResponse
	GetAllTemplates(*gin.Context) APIResponse
	GetTemplate(*gin.Context) APIResponse
	UpdateTemplate(*gin.Context) APIResponse
}

type templateHandler struct {
	TemplateService   template.Usecase
	TemplatePresenter presenter.ConvertTemplate
	Validator         validation.CustomValidator
}

func NewTemplateHandler(tuc template.Usecase, tp presenter.ConvertTemplate, v validation.CustomValidator) TemplateHandler {
	return &templateHandler{
		TemplateService:   tuc,
		TemplatePresenter: tp,
		Validator:         v,
	}
}

// CreateTemplate	godoc
// CreateTemplate	API
//
//	@Summary		Create Template
//	@Description	Create template
//	@Tags			template
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/template/create [post]
//	@Param			create_template_request	body		request.CreateTemplateRequest	true	"Create Template Request"
//	@Success		200						{object}	APIResponse{result=bool}
//	@Failure		400						{object}	APIResponse
//	@Failure		500						{object}	APIResponse
func (h *templateHandler) CreateTemplate(c *gin.Context) APIResponse {
	var request request.CreateTemplateRequest

	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	ok, code, err := h.TemplateService.CreateTemplate(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if !ok {
		return CreateResponse(err, http.StatusInternalServerError, "", common.MessageErrorCreateTemplateFail, ok)
	}
	return HandlerResponse(http.StatusOK, "", "", ok)
}

// GetTemplate	godoc
// GetTemplate	API
//
//	@Summary		Get Template
//	@Description	Get template
//	@Tags			template
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/template/{template_id} [get]
//	@Param			template_id	path		string	true	"Template ID"
//	@Success		200			{object}	APIResponse{result=presenter.TemplateWebpagesResponse}
//	@Failure		500			{object}	APIResponse
func (h *templateHandler) GetTemplate(c *gin.Context) APIResponse {
	templateID := c.Param("template_id")

	templateWebpages, code, err := h.TemplateService.GetTemplateWebpages(&templateID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.TemplatePresenter.ResponseGetTemplateWebpages(templateWebpages)
	return HandlerResponse(http.StatusOK, "", "", result)
}

// GetAllTemplates	godoc
// GetAllTemplates	API
//
//	@Summary		Get All Templates
//	@Description	get all templates
//	@Tags			template
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/template/all [get]
//	@Success		200	{object}	APIResponse{result=presenter.ListTemplateResponse}
//	@Failure		500	{object}	APIResponse
func (h *templateHandler) GetAllTemplates(c *gin.Context) APIResponse {
	templates, code, err := h.TemplateService.GetAllTemplates()
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	result := h.TemplatePresenter.ResponseGetAllTemplates(templates)
	return HandlerResponse(http.StatusOK, "", "", result)
}

// UpdateTemplate	godoc
// UpdateTemplate	API
//
//	@Summary		Update Template
//	@Description	Update Template
//	@Tags			template
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/template/{template_id} [put]
//	@Param			template_id				path		string							true	"Template ID"
//	@Param			update_template_request	body		request.UpdateTemplateRequest	true	"Update Template Request"
//	@Success		200						{object}	APIResponse{result=bool}
//	@Failure		400						{object}	APIResponse
//	@Failure		500						{object}	APIResponse
func (h *templateHandler) UpdateTemplate(c *gin.Context) APIResponse {
	var request request.UpdateTemplateRequest

	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}
	request.TemplateID = c.Param("template_id")

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}
	ok, code, err := h.TemplateService.UpdateTemplate(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if !ok {
		return CreateResponse(err, http.StatusInternalServerError, "", common.MessageErrorUpdateTemplateFail, ok)
	}
	return HandlerResponse(http.StatusOK, "", "", ok)
}

func (h *templateHandler) CloneTemplate(c *gin.Context) APIResponse {
	var req request.TemplateRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	newTemplate, code, err := h.TemplateService.CloneTemplate(&req.TemplateID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", newTemplate)
}
