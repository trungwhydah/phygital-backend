package handler

import (
	"backend-service/internal/core_backend/common/helper"
	"backend-service/internal/core_backend/entity"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/usecase/organization"
)

// OrganizationHandler interface
type OrganizationHandler interface {
	CreateOrganization(*gin.Context) APIResponse
	UpdateOrganization(*gin.Context) APIResponse
	GetAllOrganizations(*gin.Context) APIResponse
	GetOrganization(*gin.Context) APIResponse
}

// organizationHandler struct
type organizationHandler struct {
	OrganizationService   organization.UseCase
	OrganizationPresenter presenter.ConvertOrganization
	Validator             validation.CustomValidator
}

// NewOrganizationHandler create handler
func NewOrganizationHandler(ouc organization.UseCase, dp presenter.ConvertOrganization, v validation.CustomValidator) OrganizationHandler {
	return &organizationHandler{
		OrganizationService:   ouc,
		OrganizationPresenter: dp,
		Validator:             v,
	}
}

// CreateOrganization	godoc
// CreateOrganization	API
//
//	@Summary		Create Organization
//	@Description	Create organization
//	@Tags			organization
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/organization/create [post]
//	@Param			create_organization_request	formData	request.CreateOrganizationRequest	true	"create_organization_request"
//	@Success		200							{object}	APIResponse{result=presenter.OrganizationResponse}
//	@Failure		500							{object}	APIResponse
func (h *organizationHandler) CreateOrganization(c *gin.Context) APIResponse {
	//Check user role
	userRole, err := GetRoleFromGinContext(c)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}

	if userRole != string(entity.SUPER_ADMIN_ROLE) {
		err = errors.New("Unauthorized: only super admin can create organization")
		return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
	}

	// Process request
	var request request.CreateOrganizationRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	ok, err := helper.ValidateOrganizationTagName(request.NameTag)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}
	if !ok {
		err = errors.New(common.MessageErrorInvalidOrgTagName)
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	organization, code, err := h.OrganizationService.CreateOrganization(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.OrganizationPresenter.OrganizationResponse(organization)

	return HandlerResponse(code, "", "", result)
}

// GetAllOrganization	godoc
// GetAllOrganization	API
//
//	@Summary		Get All Organizations
//	@Description	Get all organizations
//	@Tags			organization
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/organization [get]
//	@Success		200	{object}	APIResponse{result=presenter.OrganizationListResponse}
//	@Failure		500	{object}	APIResponse
func (h *organizationHandler) GetAllOrganizations(c *gin.Context) APIResponse {
	// Get Org from UserInfo
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}
	info := decodeToken.(*entity.User)

	switch info.Role {
	// Org Admin
	case string(entity.ORG_ADMIN_ROLE):
		org, code, err := h.OrganizationService.GetOrgByTagName(&info.Organization)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		orgs := []entity.Organization{*org}
		result := h.OrganizationPresenter.ResponseAllOrganization(&orgs)

		return HandlerResponse(code, "", "", result)

	// Super Admin
	case string(entity.SUPER_ADMIN_ROLE):
		allOrgs, code, err := h.OrganizationService.GetAllOrganizations()
		if err != nil {
			return CreateResponse(err, http.StatusBadRequest, "", "", nil)
		}

		result := h.OrganizationPresenter.ResponseAllOrganization(allOrgs)

		return HandlerResponse(code, "", "", result)

	// Otherwise
	default:
		return CreateResponse(errors.New("Invalid user role"), http.StatusUnauthorized, "", "Invalid user role", nil)
	}
}

// UpdateOrganization	godoc
// UpdateOrganization	API
//
//	@Summary		Update Organization
//	@Description	Update Organization
//	@Tags			organization
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/organization/{org_id} [put]
//	@Param			org_id						path		string								true	"Organization ID Param"
//	@Param			update_organization_request	formData	request.UpdateOrganizationRequest	true	"Update Organization Request"
//	@Success		200							{object}	APIResponse{result=bool}
//	@Failure		500							{object}	APIResponse
func (h *organizationHandler) UpdateOrganization(c *gin.Context) APIResponse {
	// Get OrgTagName from UserInfo
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}
	info := decodeToken.(*entity.User)

	orgID := c.Param("org_id")

	// Check if org_admin has the same org_id
	switch info.Role {
	case string(entity.ORG_ADMIN_ROLE):
		org, code, err := h.OrganizationService.GetDetailOrganization(&orgID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if org.NameTag != info.Organization {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	var request request.UpdateOrganizationRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	request.OrgID = orgID

	if err := h.Validator.Validate(request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	ok, err := helper.ValidateOrganizationTagName(request.OrgTagName)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	if !ok {
		err = errors.New(common.MessageErrorInvalidOrgTagName)
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	success, code, err := h.OrganizationService.UpdateOrganization(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", success)
}

// GetOrganization	godoc
// GetOrganization	API
//
//	@Summary		Get Detail Organization
//	@Description	Get Detail Organization
//	@Tags			organization
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/organization/{org_tag_name} [get]
//	@Param			org_tag_name	path		string	true	"Organization Tag Name"
//	@Success		200				{object}	APIResponse{result=presenter.OrganizationResponse}
//	@Failure		500				{object}	APIResponse
func (h *organizationHandler) GetOrganization(c *gin.Context) APIResponse {
	userRole, err := GetRoleFromGinContext(c)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}

	orgTagName := c.Param("org_tag_name")

	//Check if org_admin has the same org_tag_name
	if userRole == string(entity.ORG_ADMIN_ROLE) {
		//Get OrgTagName from UserInfo
		decodeToken, isExisted := c.Get("userInfo")
		if !isExisted {
			return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
		}
		info := decodeToken.(*entity.User)
		if orgTagName != info.Organization {
			err = errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	org, code, err := h.OrganizationService.GetOrgByTagName(&orgTagName)

	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	if org == nil {
		return CreateResponse(
			errors.New(common.MessageErrorNotFoundOrganization),
			http.StatusBadRequest,
			"",
			common.MessageErrorNotFoundOrganization,
			nil,
		)
	}

	result := h.OrganizationPresenter.OrganizationResponse(org)

	return HandlerResponse(code, "", "", result)
}
