package presenter

import (
	"backend-service/internal/core_backend/entity"
)

// OrganizationResponse data struct
type OrganizationResponse struct {
	OrgID      string `json:"org_id"`
	OrgName    string `json:"org_name"`
	OrgTagName string `json:"org_tag_name"`
	OrgLogoURL string `json:"org_logo_url"`
}

type OrganizationListResponse struct {
	OrganizationList []OrganizationResponse `json:"org_list"`
}

// presenterOrganization struct
type PresenterOrganization struct{}

// presenterOrganization interface
type ConvertOrganization interface {
	OrganizationResponse(organization *entity.Organization) *OrganizationResponse
	ResponseAllOrganization(organization *[]entity.Organization) *OrganizationListResponse
}

// NewPresenterOrganization Constructs presenter
func NewPresenterOrganization() ConvertOrganization {
	return &PresenterOrganization{}
}

// Return property data response
func (pp *PresenterOrganization) OrganizationResponse(organization *entity.Organization) *OrganizationResponse {
	response := &OrganizationResponse{
		OrgID:      organization.ID.Hex(),
		OrgName:    organization.OrganizationName,
		OrgTagName: organization.NameTag,
		OrgLogoURL: organization.LogoURL,
	}

	return response
}

func (pp *PresenterOrganization) ResponseAllOrganization(organizationList *[]entity.Organization) *OrganizationListResponse {
	var response OrganizationListResponse
	for _, organization := range *organizationList {
		response.OrganizationList = append(response.OrganizationList, OrganizationResponse{
			OrgID:      organization.ID.Hex(),
			OrgName:    organization.OrganizationName,
			OrgTagName: organization.NameTag,
			OrgLogoURL: organization.LogoURL,
		})
	}

	return &response
}
