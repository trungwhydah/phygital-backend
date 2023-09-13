package presenter

import (
	"fmt"

	config "backend-service/config/core_backend"
)

// ScanResponse data struct
type ScanResponse struct {
	// 	UID         string `json:"uid"`
	// 	ChipID      string `json:"chip_id"`
	// 	ScanCounter int    `json:"scan_counter"`
	// 	EncMode     string `json:"enc_mode"`
	URL string `json:"url"`
}

// TapResponse data struct
type TapResponse struct {
	URL string `json:"url"`
}

// presenterScan struct
type PresenterScan struct{}

// presenterScan interface
type ConvertScan interface {
	ResponseScan(tagID *string, lang, orgTagName string) *ScanResponse
	ResponseTap(tagID *string, lang, orgTagName string) *TapResponse
}

// NewPresenterScan Constructs presenter
func NewPresenterScan() ConvertScan {
	return &PresenterScan{}
}

// Return property data response
func (pp *PresenterScan) ResponseScan(tagID *string, lang, orgTagName string) *ScanResponse {
	response := &ScanResponse{
		URL: fmt.Sprintf("%s/%s/%s/%s", config.C.Domains.WebpageDomain, lang, orgTagName, *tagID),
	}

	return response
}

func (pp *PresenterScan) ResponseTap(tagID *string, lang, orgTagName string) *TapResponse {
	response := &TapResponse{
		URL: fmt.Sprintf("%s/%s/%s/%s", config.C.Domains.WebpageDomain, lang, orgTagName, *tagID),
	}

	return response
}
