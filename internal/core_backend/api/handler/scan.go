package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/common/logger"
	validation "backend-service/internal/core_backend/infrastructure/validator"

	"github.com/gin-gonic/gin"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/mapping"
	"backend-service/internal/core_backend/usecase/organization"
	"backend-service/internal/core_backend/usecase/product"
	"backend-service/internal/core_backend/usecase/productItem"
	"backend-service/internal/core_backend/usecase/scan"
	"backend-service/internal/core_backend/usecase/session"
	"backend-service/internal/core_backend/usecase/tag"
	"backend-service/internal/core_backend/usecase/template"
	"backend-service/internal/core_backend/usecase/verification"
)

// ScanHandler interface
type ScanHandler interface {
	DecodeScan(*gin.Context) RedirectResponse
	Tap(*gin.Context) RedirectResponse
}

// scanHandler struct
type scanHandler struct {
	SessionService      session.UseCase
	OrganizationService organization.UseCase
	VerificationService verification.UseCase
	TagService          tag.UseCase
	ScanService         scan.UseCase
	MappingService      mapping.UseCase
	ProductService      product.UseCase
	ProductItemService  productItem.UseCase
	TemplateService     template.Usecase
	ScanPresenter       presenter.ConvertScan
	Validator           validation.CustomValidator
}

// NewScanHandler create handler
func NewScanHandler(ssuc session.UseCase, ouc organization.UseCase, vuc verification.UseCase, puc product.UseCase, piuc productItem.UseCase, cuc tag.UseCase, suc scan.UseCase, muc mapping.UseCase, tuc template.Usecase, dp presenter.ConvertScan, v validation.CustomValidator) ScanHandler {
	return &scanHandler{
		SessionService:      ssuc,
		OrganizationService: ouc,
		VerificationService: vuc,
		ProductService:      puc,
		ProductItemService:  piuc,
		TagService:          cuc,
		ScanService:         suc,
		MappingService:      muc,
		TemplateService:     tuc,
		ScanPresenter:       dp,
		Validator:           v,
	}
}

// DecodeScan	godoc
// DecodeScan	API
//
//	@Summary		NFC Verify Scan
//	@Description	Nfc verify scan
//	@Tags			scan
//	@Accept			multipart/form-data
//	@Produce		json
//	@Router			/scan/nfc/verify [get]
//	@Param			scan_request	formData	request.ScanRequest	true	"Scan Request"
//	@Success		302				{object}	RedirectResponse
//	@Failure		303				{object}	RedirectResponse
func (h *scanHandler) DecodeScan(c *gin.Context) RedirectResponse {
	errorPageURL := fmt.Sprintf("%s/%s", config.C.Domains.WebpageDomain, config.C.Domains.ScanErrorPage)
	var request request.ScanRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Println("[DEBUG] - 1 - Error: Binding Request: ", err.Error())
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if e := h.Validator.Validate(request); e != nil {
		log.Println("[DEBUG] - 2 - Error: Validate Request: ", e.Error())
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	scanInfo, _, err := h.ScanService.ProcessScan(&request)
	if err != nil {
		log.Println("[DEBUG] - 20 - ProcessScan")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	tag, _, err := h.TagService.GetTagByHWID(&scanInfo.UID)
	if err != nil || tag == nil {
		log.Println("[DEBUG] - 21 - GetTagByHWID")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}
	// Verification section
	verification, _, err := h.VerificationService.Verify(scanInfo, tag)
	if err != nil {
		log.Println("[DEBUG] - 22 - Verify")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	lang := "vi"
	if !verification.IsValid {
		// Temp: Check and redirect to Ortho Error page
		// TODO: Refactor to Organization section
		splitedTagID := strings.Split(tag.TagID, "-")
		if splitedTagID[0] == "0004" {
			return RedirectResponse{StatusCode: http.StatusSeeOther, URL: "https://ortho.fashion/non-genuine"}
		}

		log.Println("[DEBUG] - 23 - !verification.IsValid")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}
	h.TagService.UpdateTagCounter(&tag.TagID, &verification.Nonce)

	mapping, _, err := h.MappingService.GetMappingWithTagID(&tag.TagID)
	if err != nil || mapping == nil {

		log.Println("[DEBUG] - 24 - GetMappingWithTagID")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if len(mapping.ExternalURL) != 0 {
		session, _, err := h.SessionService.CreateSession(&tag.TagID, config.C.Server.SessionTimeoutInSecond)
		if err != nil {
			log.Println("[DEBUG] - 25 - CreateSession")
			return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
		}

		return RedirectResponse{StatusCode: http.StatusFound, URL: fmt.Sprintf("%s&sessionId=%s", mapping.ExternalURL, session.SessionID.Hex())}
	}

	if mapping.ProductItemID.IsZero() {
		log.Println("[DEBUG] - 10 - Empty Product Item")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	piID := mapping.ProductItemID.Hex()
	item, _, err := h.ProductItemService.GetDetailProductItem(&piID)
	if err != nil || item == nil || item.ProductID.IsZero() {
		log.Println("[DEBUG] - 26 - GetDetailProductItem")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	pID := item.ProductID.Hex()
	product, _, err := h.ProductService.GetProductByID(&pID)
	if err != nil || product.TemplateID.IsZero() {
		log.Println("[DEBUG] - 13 - Empty Product")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	tID := product.TemplateID.Hex()
	template, _, err := h.TemplateService.GetTemplate(&tID)
	if err != nil {

		log.Println("[DEBUG] - 27 - GetTemplate")
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if len(template.Languages) != 0 {
		lang = template.Languages[0]
	}

	oID := mapping.OrganizationID.Hex()
	org, _, err := h.OrganizationService.GetDetailOrganization(&oID)
	if err != nil {
		logger.LogError("[DEBUG] - 16 - Get error when verifying: " + err.Error())
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	result := h.ScanPresenter.ResponseScan(&tag.TagID, lang, org.NameTag)

	return RedirectResponse{StatusCode: http.StatusFound, URL: result.URL}
}

// Tap	godoc
// Tap	API
//
//	@Summary		NFC Tap
//	@Description	nfc tap
//	@Tags			scan
//	@Produce		json
//	@Router			/scan/nfc/tap/{tag_id} [get]
//	@Param			tag_id	path		string	true	"Tag ID"
//	@Success		302		{object}	RedirectResponse
//	@Failure		303		{object}	RedirectResponse
func (h *scanHandler) Tap(c *gin.Context) RedirectResponse {
	var request request.TapRequest
	request.TagID = c.Param("tag_id")

	errorPageURL := fmt.Sprintf("%s/%s", config.C.Domains.WebpageDomain, config.C.Domains.ScanErrorPage)
	if e := h.Validator.Validate(request); e != nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	tag, _, err := h.TagService.GetTag(request.TagID)
	if err != nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	verification, _, err := h.VerificationService.Verify(nil, tag)
	if err != nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if verification.IsValid {
		h.TagService.UpdateTagCounter(&tag.TagID, &verification.Nonce)
	}

	// Remove
	lang := "vi"
	tagID := tag.TagID
	mapping, _, err := h.MappingService.GetMappingWithTagID(&request.TagID)
	if err != nil || mapping == nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if len(mapping.ExternalURL) != 0 {
		session, _, err := h.SessionService.CreateSession(&tagID, config.C.Server.SessionTimeoutInSecond)
		if err != nil {
			return RedirectResponse{StatusCode: http.StatusFound, URL: fmt.Sprintf("%s/%s", config.C.Domains.WebpageDomain, config.C.Domains.ScanErrorPage)}
		}
		return RedirectResponse{StatusCode: http.StatusFound, URL: fmt.Sprintf("%s?sessionId=%s", mapping.ExternalURL, session.SessionID.Hex())}
	}

	if mapping.ProductItemID.IsZero() {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	piID := mapping.ProductItemID.Hex()
	item, _, err := h.ProductItemService.GetDetailProductItem(&piID)
	if err != nil || item == nil || item.ProductID.IsZero() {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	pID := item.ProductID.Hex()
	product, _, err := h.ProductService.GetProductByID(&pID)
	if err != nil || product.TemplateID.IsZero() {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	tID := product.TemplateID.Hex()
	template, _, err := h.TemplateService.GetTemplate(&tID)
	if err != nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	if len(template.Languages) != 0 {
		lang = template.Languages[0]
	}

	oID := mapping.OrganizationID.Hex()
	org, _, err := h.OrganizationService.GetDetailOrganization(&oID)
	if err != nil {
		return RedirectResponse{StatusCode: http.StatusSeeOther, URL: errorPageURL}
	}

	result := h.ScanPresenter.ResponseTap(&tagID, lang, org.NameTag)
	return RedirectResponse{StatusCode: http.StatusFound, URL: result.URL}
}
