package handler

import (
	"errors"
	"net/http"
	"strconv"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/digitalAsset"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/product"
	"backend-service/internal/core_backend/usecase/productItem"
	"backend-service/internal/core_backend/usecase/tag"
	"backend-service/internal/core_backend/usecase/template"
	"backend-service/internal/core_backend/usecase/user"

	"github.com/gin-gonic/gin"
)

// MappingHandler interface
type MappingHandler interface {
	UpdateMapping(*gin.Context) APIResponse
	Unmap(*gin.Context) APIResponse
	GetAllMapping(*gin.Context) APIResponse
	GetAllMappingForProduct(*gin.Context) APIResponse
	GetTemplateByTagID(*gin.Context) APIResponse
	MultipleMappingWithSingleProduct(*gin.Context) APIResponse
}

// mappingHandler struct
type mappingHandler struct {
	UserService         user.UseCase
	OrganizationService organization.UseCase
	TagService          tag.UseCase
	ProductService      product.UseCase
	ProductItemService  productItem.UseCase
	TemplateService     template.Usecase
	MappingService      mapping.UseCase
	MappingPresenter    presenter.ConvertMapping
	Validator           validation.CustomValidator
	DigitalAssetService digitalAsset.UseCase
}

// NewMappingHandler create handler
func NewMappingHandler(uuc user.UseCase, ouc organization.UseCase, cuc tag.UseCase, puc product.UseCase, piuc productItem.UseCase, tuc template.Usecase, ds mapping.UseCase, dp presenter.ConvertMapping, v validation.CustomValidator, dauc digitalAsset.UseCase) MappingHandler {
	return &mappingHandler{
		UserService:         uuc,
		OrganizationService: ouc,
		TagService:          cuc,
		ProductService:      puc,
		ProductItemService:  piuc,
		TemplateService:     tuc,
		MappingService:      ds,
		MappingPresenter:    dp,
		Validator:           v,
		DigitalAssetService: dauc,
	}
}

// UpdateMapping	godoc
// UpdateMapping	API
//
//	@Summary		Update Mapping
//	@Description	Update mapping (Identify mapping by tag_id, Only update fields that are given in json body)
//	@Tags			mapping
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/mapping/{tag_id} [put]
//	@Param			tag_id					path		string							true	"Tag ID"
//	@Param			update_mapping_request	body		request.UpdateMappingRequest	true	"Update Mapping Request"
//	@Success		200						{object}	APIResponse{result=bool}
//	@Failure		400						{object}	APIResponse
func (h *mappingHandler) UpdateMapping(c *gin.Context) APIResponse {
	tagID := c.Param("tag_id")
	var request request.UpdateMappingRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	if request.ProductItemID != nil {
		// Check if given product item is already mapped
		pItemID := request.ProductItemID.Hex()
		isMapped, code, err := h.ProductItemService.CheckProductItemMapped(&pItemID)

		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if isMapped {
			err = errors.New("Product item is already mapped with some tag")
			return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
		}

		// Check if given mapping is already mapped
		mapping, code, err := h.MappingService.GetMappingWithTagID(&tagID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if !mapping.ProductItemID.IsZero() {
			err := errors.New("Tag ID is already mapped with some product item")
			return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
		}
	}

	// Update Mapping
	ok, code, err := h.MappingService.UpdateMapping(&tagID, &request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", ok)
}

// Unmap	godoc
// Unmap	API
//
//	@Summary		Unmap Tag and Product Item (that is currently mapped)
//	@Description	Unmap Tag and Product Item (that is currently mapped)
//	@Tags			mapping
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/mapping/{tag_id} [delete]
//	@Param			tag_id			path		string					true	"Tag ID"
//	@Param			unmap_request	body		request.UnmapRequest	true	"Unmap Request"
//	@Success		200				{object}	APIResponse{result=bool}
//	@Failure		400				{object}	APIResponse
func (h *mappingHandler) Unmap(c *gin.Context) APIResponse {
	tagID := c.Param("tag_id")
	var request request.UnmapRequest
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	// Check if given mapping is already mapped
	mapping, code, err := h.MappingService.GetMappingWithTagID(&tagID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if mapping.ProductItemID.Hex() != request.ProductItemID {
		err := errors.New("The product item mapped with the given tag_id is different that the given product item")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	if !mapping.DigitalAssetID.IsZero() {
		err := errors.New("This mapping already have digital asset, couldn't unmap this")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	// Update Mapping
	ok, code, err := h.MappingService.Unmap(&tagID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", ok)
}

// GetAllMapping	godoc
// GetAllMapping	API
//
//	@Summary		Get All Mapping (Of 1 Org/All Orgs)
//	@Description	Get all mapping (of 1 org or all orgs (required super admin, leave org tag name empty))
//	@Tags			mapping
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/mapping [get]
//	@Param			org_tag_name	query		string	false	"Organization Tag Name Query"
//	@Param			is_mapped		query		bool	false	"Is Already Mapped? Query"
//	@Success		200				{object}	APIResponse{result=presenter.GetAllMappingResponse}
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *mappingHandler) GetAllMapping(c *gin.Context) APIResponse {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}

	info := decodeToken.(*entity.User)
	org_tag_name, exist := c.GetQuery("org_tag_name")
	var orgID string
	if exist {
		if info.Role == string(entity.ORG_ADMIN_ROLE) && info.Organization != org_tag_name {
			err := errors.New("Unauthorized: You can't get mappings that are NOT in your organization!")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
		org, code, err := h.OrganizationService.GetOrgByTagName(&org_tag_name)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		orgID = org.ID.Hex()
	} else {
		if info.Role != string(entity.SUPER_ADMIN_ROLE) {
			err := errors.New("Missing org tag name query, only super admin can view all mappings")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	mappingsRaw, code, err := h.MappingService.GetAllMappingInOrg(&orgID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	var mappings []entity.Mapping
	isMappedStr, exist := c.GetQuery("is_mapped")
	if exist {
		isMapped, _ := strconv.ParseBool(isMappedStr)
		for _, m := range *mappingsRaw {
			if m.ProductItemID.IsZero() == !isMapped {
				mappings = append(mappings, m)
			}
		}
	} else {
		mappings = *mappingsRaw
	}

	var owners []entity.User
	var assets []entity.DigitalAsset
	for _, m := range mappings {
		if len(m.OwnerID) != 0 {
			owner, code, err := h.UserService.GetUserByID(&m.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			owners = append(owners, *owner)
		}
		if !m.DigitalAssetID.IsZero() {
			daID := m.DigitalAssetID.Hex()
			da, code, err := h.DigitalAssetService.GetDigitalAssetByID(&daID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			assets = append(assets, *da)
		}
	}

	result := h.MappingPresenter.ResponseGetAllMapping(&mappings, &owners, &assets)

	return HandlerResponse(code, "", "", result)
}

// GetAllMappingForProduct	godoc
// GetAllMappingForProduct	API
//
//	@Summary		Get All Mapping For Product (Unmapped Mappings And Mappings Of That Product)
//	@Description	Get All Mapping For Product (Unmapped Mappings And Mappings Of That Product)
//	@Tags			mapping
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/mapping/product/{product_id} [get]
//	@Param			product_id	path		string	true	"Product ID Param"
//	@Param			is_mapped	query		bool	false	"Is Already Mapped? Query"
//	@Success		200			{object}	APIResponse{result=presenter.GetAllMappingResponse}
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
func (h *mappingHandler) GetAllMappingForProduct(c *gin.Context) APIResponse {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}

	info := decodeToken.(*entity.User)
	productID := c.Param("product_id")
	product, code, err := h.ProductService.GetProductByID(&productID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	oID := product.OrganizationID.Hex()
	switch info.Role {
	case string(entity.ORG_ADMIN_ROLE):
		org, code, err := h.OrganizationService.GetDetailOrganization(&oID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if info.Organization != org.NameTag {
			err = errors.New("Unauthorized: Given product is NOT in your organization!")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	case string(entity.SUPER_ADMIN_ROLE):
	default:
		err := errors.New("Unidentified user role in token")
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	mappingsRaw, code, err := h.MappingService.GetAllMappingForProduct(&productID, &oID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	var mappings []entity.Mapping
	isMappedStr, exist := c.GetQuery("is_mapped")
	if exist {
		isMapped, _ := strconv.ParseBool(isMappedStr)
		for _, m := range *mappingsRaw {
			if m.ProductItemID.IsZero() == !isMapped {
				mappings = append(mappings, m)
			}
		}
	} else {
		mappings = *mappingsRaw
	}
	var owners []entity.User
	var assets []entity.DigitalAsset
	for _, m := range mappings {
		if len(m.OwnerID) != 0 {
			owner, code, err := h.UserService.GetUserByID(&m.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			owners = append(owners, *owner)
		}
		if !m.DigitalAssetID.IsZero() {
			daID := m.DigitalAssetID.Hex()
			da, code, err := h.DigitalAssetService.GetDigitalAssetByID(&daID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			assets = append(assets, *da)
		}
	}

	result := h.MappingPresenter.ResponseGetAllMapping(&mappings, &owners, &assets)

	return HandlerResponse(code, "", "", result)
}

func (h *mappingHandler) GetTemplateByTagID(c *gin.Context) APIResponse {
	return HandlerResponse(200, "", "", nil)
}

// MultipleMappingWithSingleProduct	godoc
// MultipleMappingWithSingleProduct	API
//
//	@Summary		Multiple Mapping With Single Product
//	@Description	multiple mapping with single product
//	@Tags			mapping
func (h *mappingHandler) MultipleMappingWithSingleProduct(c *gin.Context) APIResponse {
	return HandlerResponse(200, "", "", nil)
}
