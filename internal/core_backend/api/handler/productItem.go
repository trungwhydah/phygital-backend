package handler

import (
	"backend-service/internal/core_backend/usecase/author"
	"errors"
	"math/big"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/entity"
	validation "backend-service/internal/core_backend/infrastructure/validator"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/digitalAsset"
	"backend-service/internal/core_backend/usecase/digitalAssetCollection"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/nft"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/product"
	"backend-service/internal/core_backend/usecase/productItem"
	"backend-service/internal/core_backend/usecase/template"
	"backend-service/internal/core_backend/usecase/user"
	webpage "backend-service/internal/core_backend/usecase/webPage"
)

// ItemHandler interface
type ProductItemHandler interface {
	CreateProductItem(*gin.Context) APIResponse
	CreateMultipleProductItems(*gin.Context) APIResponse
	GetDetailProductItem(*gin.Context) APIResponse
	GetMetadata(*gin.Context) APIResponse
	GetAllProductItem(c *gin.Context) APIResponse
	GetAllProductItemInOrg(c *gin.Context) APIResponse
	ClaimItem(*gin.Context) APIResponse
	GetStoryByTagID(*gin.Context) APIResponse
	ToggleClaimableItem(*gin.Context) APIResponse
	LikeProductItem(*gin.Context) APIResponse
	GetGalleryOfProductItemsInOrg(*gin.Context) APIResponse
	GetGalleryOfProductItemsInOrgV2(*gin.Context) APIResponse
	MintProductItem(*gin.Context) APIResponse
}

// itemHandler struct
type productItemHandler struct {
	UserService                   user.UseCase
	ProductService                product.UseCase
	ProductItemService            productItem.UseCase
	ProductItemPresenter          presenter.ConvertProductItem
	MappingService                mapping.UseCase
	OrganizationService           organization.UseCase
	TemplateService               template.Usecase
	WebpageService                webpage.UseCase
	Validator                     validation.CustomValidator
	DigitalAssetService           digitalAsset.UseCase
	DigitalAssetCollectionService digitalAssetCollection.UseCase
	NFTService                    nft.UseCase
	AuthorService                 author.UseCase
}

// NewItemHandler create handler
func NewProductItemHandler(uuc user.UseCase, puc product.UseCase, piuc productItem.UseCase, dp presenter.ConvertProductItem, m mapping.UseCase, o organization.UseCase, t template.Usecase, w webpage.UseCase, v validation.CustomValidator, duc digitalAsset.UseCase, cuc digitalAssetCollection.UseCase, nft nft.UseCase, author author.UseCase) ProductItemHandler {
	return &productItemHandler{
		UserService:                   uuc,
		ProductService:                puc,
		ProductItemService:            piuc,
		ProductItemPresenter:          dp,
		MappingService:                m,
		OrganizationService:           o,
		TemplateService:               t,
		WebpageService:                w,
		Validator:                     v,
		DigitalAssetService:           duc,
		DigitalAssetCollectionService: cuc,
		NFTService:                    nft,
		AuthorService:                 author,
	}
}

// CreateProductItem	godoc
// CreateProductItem	API
//
//	@Summary		Create Product Item
//	@Description	Create product item
//	@Tags			product-item
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product-item/create [post]
//	@Param			create_product_item_request	formData	request.CreateProductItemRequest	true	"Create Product Item Request"
//	@Success		200							{object}	APIResponse{result=presenter.ProductItemDetailResponse}
//	@Failure		400							{object}	APIResponse
func (h *productItemHandler) CreateProductItem(c *gin.Context) APIResponse {
	// TODO: Will be check owner of the product and organzation ID
	var req request.CreateProductItemRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}

	item, code, err := h.ProductItemService.CreateProductItem(&req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	var reqProduct request.InteractProductDetailRequest
	reqProduct.ProductID = item.ProductID.Hex()
	product, code, err := h.ProductService.GetProductDetail(&reqProduct)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.ProductItemPresenter.ResponseProductItemDetail(item, product)

	return HandlerResponse(code, "", "", result)
}

// CreateMultipleProductItems	godoc
// CreateMultipleProductItems	API
//
//	@Summary		Create Multiple Product Item
//	@Description	Create Multiple Product Item
//	@Tags			product-item
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product-item/create-multiple [post]
//	@Param			create_multiple_product_items_request	formData	request.CreateMultipleProductItemsRequest	true	"Create Product Item Request"
//	@Success		200										{object}	APIResponse{result=bool}
//	@Failure		400										{object}	APIResponse
func (h *productItemHandler) CreateMultipleProductItems(c *gin.Context) APIResponse {
	// TODO: Will be check owner of the product and organzation ID
	var req request.CreateMultipleProductItemsRequest
	if err := c.ShouldBind(&req); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}
	product, code, err := h.ProductService.GetProductByID(&req.ProductID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	ok, code, err := h.ProductItemService.CreateMultipleProductItems(&req.ProductID, req.NumItems, product.TotalItem+1)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	totalItems, code, err := h.ProductItemService.CountNumProductItems(&req.ProductID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	ok, code, err = h.ProductService.SyncTotalItems(&req.ProductID, totalItems)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", ok)
}

// GetDetailProductItem	godoc
// GetDetailProductItem	API
//
//	@Summary		Get Detail Of Product Item
//	@Description	Get detail of product item
//	@Tags			product-item user
//	@Produce		json
//	@Router			/product-item/{product_item_id}/detail [get]
//	@Param			product_item_id	path		string	true	"Product Item ID"
//	@Success		200				{object}	APIResponse{result=presenter.ProductItemDetailResponse}
//	@Failure		400				{object}	APIResponse
func (h *productItemHandler) GetDetailProductItem(c *gin.Context) APIResponse {
	var req request.ProductItemInteractionRequest
	req.ProductItemID = c.Param("product_item_id")
	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}

	item, code, err := h.ProductItemService.GetDetailProductItem(&req.ProductItemID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	var reqProduct request.InteractProductDetailRequest
	reqProduct.ProductID = item.ProductID.Hex()
	product, code, err := h.ProductService.GetProductDetail(&reqProduct)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	result := h.ProductItemPresenter.ResponseProductItemDetail(item, product)

	return HandlerResponse(code, "", "", result)
}

// GetAllProductItem	godoc
// GetAllProductItem	API
//
//	@Summary		Get all product items with filter
//	@Description	Return ProductItems with is_mapped and product_id filters
//	@Tags			product-item
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product-item [get]
//	@Param			product_id	query		string	true	"Product ID Query"
//	@Param			is_mapped	query		bool	false	"Is Already Mapped? Query"
//	@Success		200			{object}	APIResponse{result=presenter.AllProductItemResponse}
//	@Failure		400			{object}	APIResponse
func (h *productItemHandler) GetAllProductItem(c *gin.Context) APIResponse {
	productID, exist := c.GetQuery("product_id")
	if !exist {
		err := errors.New("Required product_id query")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	product, code, err := h.ProductService.GetProductByID(&productID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	productItems, code, err := h.ProductItemService.GetAllProductItemInProduct(&productID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	isMappedStr, exist := c.GetQuery("is_mapped")
	if exist {
		isMapped, err := strconv.ParseBool(isMappedStr)
		if err != nil {
			return CreateResponse(err, http.StatusBadRequest, "", "Query is_mapped must be of boolean type", nil)
		}
		var filteredProductItems []entity.ProductItem
		for _, productItem := range *productItems {
			pID := productItem.ID.Hex()
			isAlreadyMapped, code, err := h.MappingService.IsProductItemIDMapped(&pID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			if isAlreadyMapped == isMapped {
				filteredProductItems = append(filteredProductItems, productItem)
			}
		}
		productItems = &filteredProductItems
	}

	var users []entity.User
	for _, productItem := range *productItems {
		pID := productItem.ID.Hex()
		mapping, code, err := h.MappingService.GetMappingWithProductItemID(&pID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if mapping != nil && mapping.OwnerID != "" {
			user, code, err := h.UserService.GetUserByID(&mapping.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			users = append(users, *user)
		} else {
			users = append(users, entity.User{})
		}
	}

	result := h.ProductItemPresenter.ResponseGetAllProductItemsOfProduct(productItems, &users, &product.ProductName)

	return HandlerResponse(code, "", "", result)
}

func (h *productItemHandler) GetMetadata(c *gin.Context) APIResponse {
	tokenId := c.Param("token_id")
	tId := new(big.Int)

	tId, success := tId.SetString(tokenId, 10)
	if !success {
		err := errors.New("Require token_id")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	result := h.ProductItemPresenter.ResponseGetMetadata(tId)

	return HandlerResponse(http.StatusOK, "", "", result)
}

// GetAllProductItemInOrg	godoc
// GetAllProductItemInOrg	API
//
//	@Summary		Get All Product Items In A Specific Organization
//	@Description	Get All Product Items In A Specific Organization
//	@Tags			product-item
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product-item/organization/{org_tag_name} [get]
//	@Param			org_tag_name	path		string	true	"Organization Tag Name Param"
//	@Success		200				{object}	APIResponse{result=presenter.AllProductItemResponse}
//	@Failure		400				{object}	APIResponse
func (h *productItemHandler) GetAllProductItemInOrg(c *gin.Context) APIResponse {
	orgTagName := c.Param("org_tag_name")
	org, code, err := h.OrganizationService.GetOrgByTagName(&orgTagName)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	oID := org.ID.Hex()
	decodeToken, isExisted := c.Get("userInfo")
	if !isExisted {
		return CreateResponse(errors.New(common.MessageErrorFailDetectUser), http.StatusNonAuthoritativeInfo, "", common.MessageErrorFailDetectUser, nil)
	}

	info := decodeToken.(*entity.User)
	switch info.Role {
	case string(entity.ORG_ADMIN_ROLE):
		if info.Organization != orgTagName {
			err = errors.New("Unauthorized: You can't get product items that are NOT in your organization!")
			return CreateResponse(err, http.StatusUnauthorized, "", err.Error(), nil)
		}
	case string(entity.SUPER_ADMIN_ROLE):
	default:
		err := errors.New("Unidentified user role in token")
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	productItems, code, err := h.ProductItemService.GetProductItemsInOrg(&oID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	var users []entity.User
	for _, productItem := range *productItems {
		pID := productItem.ID.Hex()
		mapping, code, err := h.MappingService.GetMappingWithProductItemID(&pID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if mapping != nil && mapping.OwnerID != "" {
			user, code, err := h.UserService.GetUserByID(&mapping.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			users = append(users, *user)
		} else {
			users = append(users, entity.User{})
		}
	}

	var products []entity.Product
	for _, productItem := range *productItems {
		pID := productItem.ProductID.Hex()
		product, code, err := h.ProductService.GetProductByID(&pID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		products = append(products, *product)
	}

	result := h.ProductItemPresenter.ResponseGetAllProductItemsInOrg(productItems, &users, &products)

	return HandlerResponse(code, "", "", result)
}

// ClaimItem	godoc
// ClaimItem	API
//
//	@Summary		Claim Product Item
//	@Description	Claim product item
//	@Tags			product-item user
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/product-item/{product_item_id}/claim [put]
//	@Param			product_item_id	path		string	true	"Product Item ID"
//	@Success		200				{object}	APIResponse{result=bool}
//	@Failure		400				{object}	APIResponse
func (h *productItemHandler) ClaimItem(c *gin.Context) APIResponse {
	req := request.SetOwnerRequest{
		Token:         c.GetHeader("Authorization"),
		ProductItemID: c.Param("product_item_id"),
	}

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}

	user, code, err := h.UserService.TokenToUser(&req.Token)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if user == nil {
		return CreateResponse(nil, http.StatusBadRequest, "", common.MessageErrorNotFoundUser, nil)
	}

	req.OwnerID = user.ID
	result, code, err := h.ProductItemService.SetOnwerForItem(&req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", result)
}

// GetStoryByTagID	godoc
// GetStoryByTagID	API
//
//	@Summary		Get Story Detail By Tag ID
//	@Description	Get Story Detail By Tag ID
//	@Tags			product-item user
//	@Produce		json
//	@Router			/product-item/story [get]
//	@Param			tag_id	query		string	true	"Tag ID Query"
//	@Success		200		{object}	APIResponse{result=presenter.StoryDetailResponse}
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
func (h *productItemHandler) GetStoryByTagID(c *gin.Context) APIResponse {
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

	var (
		product     *entity.Product
		productItem *entity.ProductItem
		owner       *entity.User
		homepage    *entity.WebPage
		template    *entity.TemplateWebpages
		da          *entity.DigitalAsset
		dac         *entity.DigitalAssetCollection
		at          *entity.Author
	)
	if !mapping.ProductItemID.IsZero() {
		productItemID := mapping.ProductItemID.Hex()
		productItem, code, err = h.ProductItemService.GetDetailProductItem(&productItemID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		productID := productItem.ProductID.Hex()
		product, code, err = h.ProductService.GetProductByID(&productID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	}

	if mapping.OwnerID != "" {
		owner, code, err = h.UserService.GetUserByID(&mapping.OwnerID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	}

	if !mapping.DigitalAssetID.IsZero() {
		daID := mapping.DigitalAssetID.Hex()
		da, code, err = h.DigitalAssetService.GetDigitalAssetByID(&daID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		dacID := da.CollectionID.Hex()
		dac, code, err = h.DigitalAssetCollectionService.GetCollectionByID(&dacID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	}

	if !product.TemplateID.IsZero() {
		templateID := product.TemplateID.Hex()
		template, code, err = h.TemplateService.GetTemplateWebpages(&templateID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	}

	for _, page := range template.Pages {
		if page.Type == "home" {
			homepage = &page
			break
		}
	}
	if homepage == nil {
		err = errors.New("Template doesn't have 'home' in pages (template_id = " + template.ID.Hex() + " )")
		return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
	}

	if !product.AuthorID.IsZero() {
		aID := product.AuthorID.Hex()
		at, code, err = h.AuthorService.GetAuthorDetail(&aID)
		if err != nil {
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
	}

	orgID := mapping.OrganizationID.Hex()
	organization, code, err := h.OrganizationService.GetDetailOrganization(&orgID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	result := h.ProductItemPresenter.ResponseGetStoryDetail(mapping, product, productItem, owner, template, homepage, organization, da, dac, at)

	return HandlerResponse(http.StatusOK, "", "", result)
}

// ToggleClaimableItem	godoc
// ToggleClaimableItem	API
//
//	@Summary		Toggle Claimable Item
//	@Description	Toggle Claimable Item
//	@Tags			product-item user
//	@Produce		json
//	@Router			/product-item/{product_item_id}/toggle-claimable [put]
//	@Param			product_item_id	path		string	true	"Product Item ID"
//	@Success		200				{object}	APIResponse{result=bool}
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *productItemHandler) ToggleClaimableItem(c *gin.Context) APIResponse {
	var req request.ProductItemInteractionRequest
	req.ProductItemID = c.Param("product_item_id")
	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}

	success, code, err := h.ProductItemService.ToggleClaimable(&req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", success)
}

// LikeProductItem	godoc
// LikeProductItem	API
//
//	@Summary		Like Product Item
//	@Description	Like Product Item
//	@Tags			product-item user
//	@Produce		json
//	@Router			/product-item/{product_item_id}/like [post]
//	@Param			product_item_id	path		string	true	"Product Item ID"
//	@Success		200				{object}	APIResponse{result=bool}
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *productItemHandler) LikeProductItem(c *gin.Context) APIResponse {
	var req request.ProductItemLikeRequest
	req.ProductItemID = c.Param("product_item_id")

	if e := h.Validator.Validate(req); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", e.Error(), nil)
	}

	ok, code, err := h.ProductItemService.UpdateTotalLike(&req.ProductItemID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", ok)
}

// GetGalleryOfProductItemsInOrg	godoc
// GetGalleryOfProductItemsInOrg	API
//
//	@Summary		Get Gallery Of Product Items Of A Competition
//	@Description	Get Gallery Of Product Items Of A Competition
//	@Tags			competition
//	@Produce		json
//	@Router			/competition/{org_tag_name} [get]
//	@Param			org_tag_name	path		string	true	"Organization Tag Name (Competition Name)"
//	@Success		200				{array}		APIResponse{result=presenter.GalleryProductItemsListResponse}
//	@Failure		203				{object}	APIResponse
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *productItemHandler) GetGalleryOfProductItemsInOrg(c *gin.Context) APIResponse {
	org_tag_name := c.Param("org_tag_name")
	org, code, err := h.OrganizationService.GetOrgByTagName(&org_tag_name)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if org == nil {
		err = errors.New("Couldn't find organization with given org tag name")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	oID := org.ID.Hex()
	mappingsRaw, code, err := h.MappingService.GetAllMappingInOrg(&oID)
	var (
		mappings              []entity.Mapping
		products              []entity.Product
		homepages, craftsmens []entity.WebPage
		templates             []entity.TemplateWebpages
		totalLikes            []int
		das                   []entity.DigitalAsset
		dacs                  []entity.DigitalAssetCollection
	)
	for _, mapping := range *mappingsRaw {
		// TODO
		g := errgroup.Group{}
		g.Go(func() error {

			return nil
		})
		//End TODO
		if mapping.ProductItemID.IsZero() {
			continue
		}
		piID := mapping.ProductItemID.Hex()
		productItem, code, err := h.ProductItemService.GetDetailProductItem(&piID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		pID := productItem.ProductID.Hex()
		product, code, err := h.ProductService.GetProductByID(&pID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if product.TemplateID.IsZero() {
			continue
		}

		tID := product.TemplateID.Hex()
		template, code, err := h.TemplateService.GetTemplateWebpages(&tID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		var homepage, craftsmen entity.WebPage
		for _, page := range template.Pages {
			if page.Type == "home" {
				homepage = page
			}
			if page.Type == "craftsmen" {
				craftsmen = page
			}
		}
		if homepage.ID.IsZero() {
			err = errors.New("Template doesn't have 'home' in pages (template_id = " + template.ID.Hex() + " )")
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}

		if !craftsmen.ID.IsZero() {
			craftsmens = append(craftsmens, craftsmen)
		} else {
			craftsmens = append(craftsmens, entity.WebPage{})
		}

		if !mapping.DigitalAssetID.IsZero() {
			daID := mapping.DigitalAssetID.Hex()
			da, code, err := h.DigitalAssetService.GetDigitalAssetByID(&daID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			dacID := da.CollectionID.Hex()
			dac, code, err := h.DigitalAssetCollectionService.GetCollectionByID(&dacID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			das = append(das, *da)
			dacs = append(dacs, *dac)
		} else {
			das = append(das, entity.DigitalAsset{})
			dacs = append(dacs, entity.DigitalAssetCollection{})
		}

		mappings = append(mappings, mapping)
		products = append(products, *product)
		totalLikes = append(totalLikes, productItem.TotalLike)
		homepages = append(homepages, homepage)
		templates = append(templates, *template)
	}
	result := h.ProductItemPresenter.ResponseGalleryProductItems(&org.OrganizationName, totalLikes, &mappings, &products, &homepages, &craftsmens, &templates, &das, &dacs)
	return HandlerResponse(code, "", "", result)
}

// GetGalleryOfProductItemsInOrgV2	godoc
// GetGalleryOfProductItemsInOrgV2	API
//
//	@Summary		Get Gallery Of Product Items Of A Competition V2
//	@Description	Get Gallery Of Product Items Of A Competition V2
//	@Tags			competition
//	@Produce		json
//	@Router			/competition/v2/{org_tag_name} [get]
//	@Param			org_tag_name	path		string	true	"Organization Tag Name (Competition Name)"
//	@Success		200				{array}		APIResponse{result=presenter.GalleryProductItemsListResponse}
//	@Failure		203				{object}	APIResponse
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *productItemHandler) GetGalleryOfProductItemsInOrgV2(c *gin.Context) APIResponse {
	org_tag_name := c.Param("org_tag_name")
	org, code, err := h.OrganizationService.GetOrgByTagName(&org_tag_name)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if org == nil {
		err = errors.New("Couldn't find organization with given org tag name")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	oID := org.ID.Hex()
	mappingsRaw, code, err := h.MappingService.GetAllMappingInOrg(&oID)
	var (
		mappings     []*entity.Mapping
		products     []*entity.Product
		productItems []*entity.ProductItem
		owners       []*entity.User
		templates    []*entity.TemplateWebpages
		das          []*entity.DigitalAsset
		dacs         []*entity.DigitalAssetCollection
		orgs         []*entity.Organization
	)
	for _, mapping := range *mappingsRaw {
		if mapping.ProductItemID.IsZero() {
			continue
		}
		var (
			product     *entity.Product
			productItem *entity.ProductItem
			owner       *entity.User
			template    *entity.TemplateWebpages
			da          *entity.DigitalAsset
			dac         *entity.DigitalAssetCollection
			org         *entity.Organization
		)
		piID := mapping.ProductItemID.Hex()
		productItem, code, err := h.ProductItemService.GetDetailProductItem(&piID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		pID := productItem.ProductID.Hex()
		product, code, err = h.ProductService.GetProductByID(&pID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if !product.TemplateID.IsZero() {
			tID := product.TemplateID.Hex()
			template, code, err = h.TemplateService.GetTemplateWebpages(&tID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
		}

		if mapping.OwnerID != "" {
			owner, code, err = h.UserService.GetUserByID(&mapping.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
		}

		if !mapping.DigitalAssetID.IsZero() {
			daID := mapping.DigitalAssetID.Hex()
			da, code, err = h.DigitalAssetService.GetDigitalAssetByID(&daID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			dacID := da.CollectionID.Hex()
			dac, code, err = h.DigitalAssetCollectionService.GetCollectionByID(&dacID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
		}
		if !mapping.OrganizationID.IsZero() {

			orgID := mapping.OrganizationID.Hex()
			org, code, err = h.OrganizationService.GetDetailOrganization(&orgID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
		}
		mappings = append(mappings, &mapping)
		products = append(products, product)
		productItems = append(productItems, productItem)
		owners = append(owners, owner)
		templates = append(templates, template)
		das = append(das, da)
		dacs = append(dacs, dac)
		orgs = append(orgs, org)
	}
	result := h.ProductItemPresenter.ResponseGalleryProductItemsV2(&mappings, &products, &productItems, &owners, &templates, &orgs, &das, &dacs)
	return HandlerResponse(code, "", "", result)
}

// MintProductItem	godoc
// MintProductItem	API
//
//	@Summary		Mint Product Item
//	@Description	Mint Product Item
//	@Tags			product-item
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/product-item/{product_item_id}/mint [post]
//	@Param			product_item_id	path		string	true	"Product Item ID"
//	@Success		200				{object}	APIResponse{result=bool}
//	@Failure		400				{object}	APIResponse
func (h *productItemHandler) MintProductItem(c *gin.Context) APIResponse {
	pItemID := c.Param("product_item_id")
	mapping, code, err := h.MappingService.GetMappingWithProductItemID(&pItemID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	if mapping == nil {
		err := errors.New("Mapping not found with the given product item")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	if mapping.OwnerID == "" {
		err := errors.New("OwnerID not found in mapping with given product item. Only mint when there is a valid OwnerID.")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	owner, code, err := h.UserService.GetUserByID(&mapping.OwnerID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	orgID := mapping.OrganizationID.Hex()
	collection, code, err := h.DigitalAssetCollectionService.GetCollectionByOrgID(&orgID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	contractAddr := collection.ContractAddress
	txHash, code, err := h.NFTService.Mint(&contractAddr, &owner.WalletAddress)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	pItemProductOrgAggregate, code, err := h.ProductItemService.GetProductItemProductOrgAggregate(&pItemID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}
	metadata := h.DigitalAssetService.ConstructMetadata(pItemProductOrgAggregate.ItemIndex, &pItemProductOrgAggregate.OrgTagName, &pItemProductOrgAggregate.Product)
	da := &entity.DigitalAsset{
		CollectionID: collection.ID,
		BaseModel: entity.BaseModel{
			Status: "Pending",
		},
		TxHash:   txHash,
		Metadata: *metadata,
	}
	da.SetTime()
	da, code, err = h.DigitalAssetService.CreateDigitalAsset(da)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	mapping.DigitalAssetID = da.ID
	req := request.UpdateMappingRequest{
		DigitalAssetID: &da.ID,
	}
	ok, code, err := h.MappingService.UpdateMapping(&mapping.TagID, &req)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	go h.NFTService.WatchTransaction(&txHash)

	return HandlerResponse(code, "", "", ok)
}
