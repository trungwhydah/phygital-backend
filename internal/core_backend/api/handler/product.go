package handler

import (
	"errors"
	"net/http"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"
	validation "backend-service/internal/core_backend/infrastructure/validator"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/product"
	"backend-service/internal/core_backend/usecase/productItem"
	"backend-service/internal/core_backend/usecase/template"
	webpage "backend-service/internal/core_backend/usecase/webPage"
)

// ProductHandler interface
type ProductHandler interface {
	CreateProduct(*gin.Context) APIResponse
	GetAllProducts(*gin.Context) APIResponse
	GetProductDetail(c *gin.Context) APIResponse
	UpdateProductDetail(c *gin.Context) APIResponse
	DeteleProductByID(c *gin.Context) APIResponse
	GetProductByTagID(c *gin.Context) APIResponse
}

// productHandler struct
type productHandler struct {
	MappingService      mapping.UseCase
	OrganizationService organization.UseCase
	WebpageService      webpage.UseCase
	ProductService      product.UseCase
	ProductItemService  productItem.UseCase
	TemplateService     template.Usecase
	ProductPresenter    presenter.ConvertProduct
	Validator           validation.CustomValidator
}

// NewProductHandler create handler
func NewProductHandler(muc mapping.UseCase, ouc organization.UseCase, wb webpage.UseCase, ds product.UseCase, piuc productItem.UseCase, ts template.Usecase, dp presenter.ConvertProduct, v validation.CustomValidator) ProductHandler {
	return &productHandler{
		MappingService:      muc,
		OrganizationService: ouc,
		WebpageService:      wb,
		ProductService:      ds,
		ProductItemService:  piuc,
		TemplateService:     ts,
		ProductPresenter:    dp,
		Validator:           v,
	}
}

// CreateProduct	godoc
// CreateProduct	API
//
//	@Summary		Create Product
//	@Description	Create product
//	@Tags			product
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product/create [post]
//	@Param			create_product_request	body		entity.Product	true	"Create Product Request"
//	@Success		200						{object}	APIResponse{result=string}
//	@Failure		400						{object}	APIResponse
//	@Failure		500						{object}	APIResponse
func (h *productHandler) CreateProduct(c *gin.Context) APIResponse {
	var request entity.Product
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}
	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	request.SetTime()
	request.SetStatus(common.StatusActive)

	oIDHex := request.OrganizationID.Hex()
	org, code, err := h.OrganizationService.GetDetailOrganization(&oIDHex)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if org == nil {
		err = errors.New("Cannot get organization entity from the given org id")
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}
	info := decodeToken.(*entity.User)
	switch info.Role {
	case string(entity.ORG_ADMIN_ROLE):
		if org.NameTag != info.Organization {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}
	product, code, err := h.ProductService.CreateProduct(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	pID := product.ID.Hex()
	totalItems := request.TotalItem
	ok, code, err := h.ProductItemService.CreateMultipleProductItems(&pID, totalItems, 1)
	if err != nil || !ok {
		err = errors.New("Cannot create multiple product items with given total item")
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.ProductPresenter.ResponseCreateProduct(&pID)

	return HandlerResponse(code, "", "", result)
}

// GetProductDetail	godoc
// GetProductDetail	API
//
//	@Summary		Get Product Detail
//	@Description	Get product detail
//	@Tags			product
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product/{product_id} [get]
//	@Param			product_id	path		string	true	"Product ID"
//	@Success		200			{object}	APIResponse{result=presenter.ProductResponse}
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
func (h *productHandler) GetProductDetail(c *gin.Context) APIResponse {
	// TODO: Check organization of user correspond to mapping or not
	var request request.InteractProductDetailRequest
	request.ProductID = c.Param("product_id")
	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	product, code, err := h.ProductService.GetProductDetail(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	orgIDOfProduct := product.OrganizationID.Hex()
	organization, code, err := h.OrganizationService.GetDetailOrganization(&orgIDOfProduct)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if organization == nil {
		err = errors.New("Cannot get organization entity from product")
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}

	role, orgTagName, err := GetUserRoleAndOrgTagNameFromGinContext(c)
	if err != nil {
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}
	if role == string(entity.ORG_ADMIN_ROLE) {
		if orgTagName != organization.NameTag {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	result := h.ProductPresenter.ResponseGetProductDetail(product, organization)
	return HandlerResponse(code, "", "", result)
}

// GetAllProducts	godoc
// GetAllProducts	API
//
//	@Summary		Get All Products
//	@Description	Get All Products
//	@Tags			product
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product [get]
//	@Param			org_tag_name	query		string	false	"Organization Tag Name (get all products if empty)"
//	@Success		200				{array}		APIResponse{result=presenter.ListProductResponse}
//	@Failure		203				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *productHandler) GetAllProducts(c *gin.Context) APIResponse {
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}
	organizationTagName := c.Query("org_tag_name")
	info := decodeToken.(*entity.User)
	if organizationTagName == "" {
		if info.Role == string(entity.ORG_ADMIN_ROLE) {
			organizationTagName = info.Organization
		}
	}
	var orgID string
	if info.Role == string(entity.ORG_ADMIN_ROLE) {
		if organizationTagName != info.Organization {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}
	if organizationTagName != "" {
		org, code, err := h.OrganizationService.GetOrgByTagName(&organizationTagName)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if org == nil {
			err = errors.New("Invalid org tag name")
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		// TODO: Check if organization of user is the same with the org of the queried product
		orgID = org.ID.Hex()
	}

	products, code, err := h.ProductService.GetAllProductsInOrg(&orgID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	var organizations []entity.Organization
	for _, product := range *products {
		orgIDOfProduct := product.OrganizationID.Hex()
		organization, code, err := h.OrganizationService.GetDetailOrganization(&orgIDOfProduct)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if organization == nil {
			err = errors.New("Cannot get organization entity from product")
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
		organizations = append(organizations, *organization)
	}

	result := h.ProductPresenter.ResponseAllProducts(products, &organizations)

	return HandlerResponse(code, "", "", result)
}

// UpdateProductDetail	godoc
// UpdateProductDetail	API
//
//	@Summary		Update Product Detail
//	@Description	Update product detail
//	@Tags			product
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product/{product_id} [put]
//	@Param			product_id				path		string			true	"Product ID"
//	@Param			update_product_detail	body		entity.Product	true	"Update Product Request"
//	@Success		200						{object}	APIResponse{result=bool}
//	@Failure		400						{object}	APIResponse
//	@Failure		500						{object}	APIResponse
func (h *productHandler) UpdateProductDetail(c *gin.Context) APIResponse {
	var request entity.Product
	if err := c.ShouldBind(&request); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", "", nil)
	}
	productID := c.Param("product_id")

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	role, orgTagName, err := GetUserRoleAndOrgTagNameFromGinContext(c)
	if role == string(entity.ORG_ADMIN_ROLE) {
		product, code, err := h.ProductService.GetProductByID(&productID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		orgIDOfProduct := product.OrganizationID.Hex()
		organization, code, err := h.OrganizationService.GetDetailOrganization(&orgIDOfProduct)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if organization == nil {
			err = errors.New("Cannot get organization entity from product")
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
		if orgTagName != organization.NameTag {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	//Check updated Template exists
	if request.TemplateID.Hex() != "" {
		tID := request.TemplateID.Hex()
		isExisted, code, err := h.TemplateService.CheckExistedTemplate(&tID)

		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if !isExisted {
			return CreateResponse(errors.New(common.MessageErrorInvalidTemplateID), http.StatusBadRequest, "", common.MessageErrorInvalidTemplateID, nil)
		}
	}

	success, code, err := h.ProductService.UpdateProductDetail(&request, &productID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", success)
}

// DeteleProductByID	godoc
// DeteleProductByID	API
//
//	@Summary		Delete Product By ID
//	@Description	Delete product by id
//	@Tags			product
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product/{product_id} [delete]
//	@Param			product_id	path		string	true	"Product ID"
//	@Success		200			{object}	APIResponse{result=bool}
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
func (h *productHandler) DeteleProductByID(c *gin.Context) APIResponse {
	var request request.InteractProductDetailRequest
	request.ProductID = c.Param("product_id")
	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	role, orgTagName, err := GetUserRoleAndOrgTagNameFromGinContext(c)
	if role == string(entity.ORG_ADMIN_ROLE) {
		product, code, err := h.ProductService.GetProductByID(&request.ProductID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		orgIDOfProduct := product.OrganizationID.Hex()
		organization, code, err := h.OrganizationService.GetDetailOrganization(&orgIDOfProduct)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if organization == nil {
			err = errors.New("Cannot get organization entity from product")
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}

		if err != nil {
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
		if orgTagName != organization.NameTag {
			err := errors.New("Unauthorized: You do not have access to this organization")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	}

	success, code, err := h.ProductService.DeteleProductByID(&request)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", success)
}

// GetProductByTagID	godoc
// GetProductByTagID	API
//
//	@Summary		Get Product Detail By Tag ID
//	@Description	Get Product Detail By Tag ID
//	@Tags			seo
//	@Produce		json
//	@Router			/product/seo [get]
//	@Param			tag_id	query		string	true	"Tag ID Query"
//	@Success		200		{object}	APIResponse{result=entity.Product}
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
func (h *productHandler) GetProductByTagID(c *gin.Context) APIResponse {
	tag_id, exist := c.GetQuery("tag_id")
	if !exist {
		err := errors.New("Require tag_id query")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	mapping, code, err := h.MappingService.GetMappingWithTagID(&tag_id)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if mapping == nil {
		err = errors.New("No mapping found with tag_id " + tag_id)
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	if mapping.ProductItemID.IsZero() {
		err = errors.New("No product item found in mapping of tag_id " + tag_id)
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	productItemID := mapping.ProductItemID.Hex()
	productItem, code, err := h.ProductItemService.GetDetailProductItem(&productItemID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	productID := productItem.ProductID.Hex()
	product, code, err := h.ProductService.GetProductByID(&productID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(http.StatusOK, "", "", product)
}
