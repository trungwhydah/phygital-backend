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
	"backend-service/internal/core_backend/usecase/digitalAsset"
	"backend-service/internal/core_backend/usecase/digitalAssetCollection"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/product"
	"backend-service/internal/core_backend/usecase/productItem"
	"backend-service/internal/core_backend/usecase/template"
	"backend-service/internal/core_backend/usecase/user"
)

// DigitalAssetHandler interface
type DigitalAssetHandler interface {
	GetDigitalAssetsByCollection(*gin.Context) APIResponse
	GetDigitalAssets(*gin.Context) APIResponse
	GetDigitalMetadataWithID(*gin.Context) APIResponse
	SyncDigitalAssetsMetadata(*gin.Context) APIResponse
}

// digitalAssetHandler struct
type digitalAssetHandler struct {
	DigitalAssetService           digitalAsset.UseCase
	DigitalAssetCollectionService digitalAssetCollection.UseCase
	MappingService                mapping.UseCase
	UserService                   user.UseCase
	ProductItemService            productItem.UseCase
	ProductService                product.UseCase
	OrganizationService           organization.UseCase
	TemplateService               template.Usecase
	DigitalAssetPresenter         presenter.ConvertDigitalAsset
	Validator                     validation.CustomValidator
}

// NewDigitalAssetHandler create handler
func NewDigitalAssetHandler(dcs digitalAssetCollection.UseCase, ds digitalAsset.UseCase, ms mapping.UseCase, us user.UseCase, pis productItem.UseCase, ps product.UseCase, os organization.UseCase, ts template.Usecase, dp presenter.ConvertDigitalAsset, v validation.CustomValidator) DigitalAssetHandler {
	return &digitalAssetHandler{
		DigitalAssetService:           ds,
		DigitalAssetCollectionService: dcs,
		MappingService:                ms,
		UserService:                   us,
		ProductItemService:            pis,
		ProductService:                ps,
		OrganizationService:           os,
		TemplateService:               ts,
		DigitalAssetPresenter:         dp,
		Validator:                     v,
	}
}

func (h *digitalAssetHandler) GetDigitalAssetsByCollection(c *gin.Context) APIResponse {
	var request = request.GetAssetByCollectionRequest{
		CollectionID: c.Param("collectionID"),
	}
	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	digitalAssets, code, err := h.DigitalAssetService.GetDigitalAssetByCollection(&request.CollectionID)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	result := h.DigitalAssetPresenter.ResponseDigitalAssets(digitalAssets)

	return APIResponse{Code: code, Result: result}
}

// GetDigitalMetadataWithID	godoc
// GetDigitalMetadataWithID	API
//
//	@Summary		Get Digital Asset Metadata By Org tag name and token ID
//	@Description	Get Digital Asset Metadata By Org tag name and token ID
//	@Tags			digital-asset
//	@Produce		json
//	@Router			/digital-asset/{org_tag_name}/{token_id} [get]
//	@Param			org_tag_name	path		string	true	"organization tag name"
//	@Param			token_id		path		string	true	"Token ID"
//	@Success		200				{object}	APIResponse{result=entity.Metadata}
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *digitalAssetHandler) GetDigitalMetadataWithID(c *gin.Context) APIResponse {
	var request = request.GetDigitalMetadataWithIDRequest{
		OrgTagName: c.Param("org_tag_name"),
		TokenID:    c.Param("token_id"),
	}

	if e := h.Validator.Validate(request); e != nil {
		return CreateResponse(e, http.StatusBadRequest, "", "", nil)
	}

	// Get Org by tag name
	org, code, err := h.OrganizationService.GetOrgByTagName(&request.OrgTagName)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	if org == nil {
		return CreateResponse(errors.New(common.MessageErrorOrgNotFound), code, "", common.MessageErrorOrgNotFound, nil)
	}

	// Get Collection by Org
	oID := org.ID.Hex()
	collection, code, err := h.DigitalAssetCollectionService.GetCollectionByOrgID(&oID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	// Get Digital Asset by Token ID and Collection
	collectionID := collection.ID.Hex()
	tokenID := request.TokenIDToInt()
	digitalAsset, code, err := h.DigitalAssetService.GetDigitalAssetByTokenID(&collectionID, &tokenID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	if digitalAsset == nil {
		err := errors.New("Couldn't find corresponding digital asset (given token_id might not exist)")
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	return APIResponse{Code: code, Result: digitalAsset.Metadata}
}

// GetDigitalAssets	godoc
// GetDigitalAssets	API
//
//	@Summary		Get All Digital Assets Or By Org Tag Name
//	@Description	Get All Digital Assets Or By Org Tag Name
//	@Tags			digital-asset
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/digital-asset [get]
//	@Param			org_tag_name	query		string	false	"Organization Tag Name Query"
//	@Success		200				{object}	APIResponse{result=presenter.ListDigitalAssetsResponse}
//	@Failure		400				{object}	APIResponse
//	@Failure		500				{object}	APIResponse
func (h *digitalAssetHandler) GetDigitalAssets(c *gin.Context) APIResponse {
	org_tag_name, exist := c.GetQuery("org_tag_name")
	var digitalAssets *[]entity.DigitalAsset
	if exist {
		// Get Org by tag name
		org, code, err := h.OrganizationService.GetOrgByTagName(&org_tag_name)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		// Get Collection by Org
		oID := org.ID.Hex()
		collection, code, err := h.DigitalAssetCollectionService.GetCollectionByOrgID(&oID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		cID := collection.ID.Hex()
		digitalAssets, code, err = h.DigitalAssetService.GetActiveDigitalAssetByCollectionID(&cID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	} else {
		var (
			code int
			err  error
		)
		digitalAssets, code, err = h.DigitalAssetService.GetAllActiveDigitalAssets()
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
	}

	var (
		collections []entity.DigitalAssetCollection
		owners      []entity.User
	)
	for _, digitalAsset := range *digitalAssets {
		cID := digitalAsset.CollectionID.Hex()
		collection, code, err := h.DigitalAssetCollectionService.GetCollectionByID(&cID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		collections = append(collections, *collection)
		dID := digitalAsset.ID.Hex()
		mapping, code, err := h.MappingService.GetMappingByDigitalAsset(&dID)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}

		if mapping.OwnerID != "" {
			user, code, err := h.UserService.GetUserByID(&mapping.OwnerID)
			if err != nil {
				return CreateResponse(err, code, "", err.Error(), nil)
			}
			owners = append(owners, *user)
		} else {
			owners = append(owners, entity.User{})
		}

	}
	result := h.DigitalAssetPresenter.ResponseGetDetailDigitalAssets(digitalAssets, &collections, &owners)

	return APIResponse{Code: http.StatusOK, Result: result}
}

// SyncDigitalAssetsMetadata	godoc
// SyncDigitalAssetsMetadata	API
//
//	@Summary		Sync Digital Assets Metadata
//	@Description	Sync Digital Assets Metadata (As metadata format can change so use this API to update the correct one to the DB)
//	@Tags			digital-asset
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/digital-asset/sync-metadata [put]
//	@Success		200	{object}	APIResponse{result=bool}
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
func (h *digitalAssetHandler) SyncDigitalAssetsMetadata(c *gin.Context) APIResponse {
	aggregations, code, err := h.DigitalAssetService.GetDigitalAssetsProductAggregate()
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	for _, aggregation := range *aggregations {
		metadata := h.DigitalAssetService.ConstructMetadata(aggregation.ItemIndex, &aggregation.OrgTagName, &aggregation.Product)
		daID := aggregation.ID.Hex()
		ok, code, err := h.DigitalAssetService.UpdateDigitalAssetMetadata(&daID, metadata)
		if err != nil {
			return CreateResponse(err, code, "", err.Error(), nil)
		}
		if !ok {
			err = errors.New("Error update digital asset metadata")
			return CreateResponse(err, http.StatusInternalServerError, "", err.Error(), nil)
		}
	}

	return APIResponse{Code: http.StatusOK, Result: true}
}
