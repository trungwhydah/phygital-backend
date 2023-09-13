package organization

import (
	"errors"
	"net/http"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service struct
type Service struct {
	repo Repository
}

// NewService create service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// GetOrganization create organization info
func (s *Service) CreateOrganization(request *request.CreateOrganizationRequest) (*entity.Organization, int, error) {
	//Check if org Tag Name is already taken
	dupOrg, err := s.repo.GetOrgByTagName(&request.NameTag)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if dupOrg != nil {
		return nil, http.StatusBadRequest, errors.New(common.MessageErrorCreateOrgFail)
	}

	org := &entity.Organization{
		OrganizationName: request.OrgName,
		NameTag:          request.NameTag,
		BaseModel: entity.BaseModel{
			Status: common.StatusActive,
		},
	}
	if len(request.LogoURL) != 0 {
		org.LogoURL = request.LogoURL
	}

	org.SetTime()

	organization, err := s.repo.CreateOrganization(org)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return organization, http.StatusOK, nil
}

func (s *Service) GetAllOrganizations() (*[]entity.Organization, int, error) {
	orgs, err := s.repo.GetAllOrganizations()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return orgs, http.StatusOK, nil
}

// UpdateOrganization update organization info
func (s *Service) UpdateOrganization(request *request.UpdateOrganizationRequest) (bool, int, error) {
	oID, err := primitive.ObjectIDFromHex(request.OrgID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	org := &entity.Organization{
		OrganizationName: request.OrgName,
		NameTag:          request.OrgTagName,
		LogoURL:          request.OrgLogoURL,
	}
	org.ID = oID
	org.SetTime()

	success, err := s.repo.UpdateOrganization(org)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	return success, http.StatusOK, nil
}

func (s *Service) GetDetailOrganization(orgID *string) (*entity.Organization, int, error) {
	org, err := s.repo.GetDetailOrganization(orgID)
	if err != nil {
		logger.LogError("[DEBUG] - 15 - error when getting organization " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return org, http.StatusOK, nil
}

func (s *Service) GetOrgByTagName(tagName *string) (*entity.Organization, int, error) {
	org, err := s.repo.GetOrgByTagName(tagName)
	if err != nil {
		logger.LogError("error when getting organization " + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return org, http.StatusOK, nil
}
