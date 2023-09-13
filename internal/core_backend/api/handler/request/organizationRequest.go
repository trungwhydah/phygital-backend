package request

type CreateOrganizationRequest struct {
	OrgName string `form:"org_name" validate:"required"`
	NameTag string `form:"org_tag_name" validate:"required"`
	LogoURL string `form:"org_logo_url"`
}

type UpdateOrganizationRequest struct {
	OrgID      string `validate:"required" swaggerignore:"true"`
	OrgName    string `form:"org_name"`
	OrgTagName string `form:"org_tag_name"`
	OrgLogoURL string `form:"org_logo_url"`
}
