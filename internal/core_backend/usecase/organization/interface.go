package organization

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/entity"
)

// Organization interface
type Organization interface {
	// Interface for repository
	CreateOrganization(*entity.Organization) (*entity.Organization, error)
	UpdateOrganization(*entity.Organization) (bool, error)
	GetAllOrganizations() (*[]entity.Organization, error)
	GetDetailOrganization(orgID *string) (*entity.Organization, error)
	GetOrgByTagName(*string) (*entity.Organization, error)
}

// Repository interface
type Repository interface {
	Organization
}

// UseCase interface
type UseCase interface {
	// Interface for usecase - service
	CreateOrganization(organizationRquest *request.CreateOrganizationRequest) (*entity.Organization, int, error)
	UpdateOrganization(organizationRquest *request.UpdateOrganizationRequest) (bool, int, error)
	GetAllOrganizations() (*[]entity.Organization, int, error)
	GetDetailOrganization(orgID *string) (*entity.Organization, int, error)
	GetOrgByTagName(tagName *string) (*entity.Organization, int, error)
}
