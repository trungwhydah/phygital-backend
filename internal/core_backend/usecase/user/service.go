package user

import (
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common"
	"backend-service/internal/core_backend/common/helper"
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"backend-service/internal/core_backend/usecase/organization"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service struct
type Service struct {
	repo     Repository
	firebase *firebase.FirebaseClient
	orgRepo  organization.Repository
}

// NewService create service
func NewService(fb *firebase.FirebaseClient, r Repository, or organization.Repository) *Service {
	return &Service{
		repo:     r,
		firebase: fb,
		orgRepo:  or,
	}
}

// CreateUser
func (s *Service) RegisterUser(request *request.CreateUserRequest) (*entity.User, int, error) {
	tokenStringArr := strings.Split(request.Token, " ")
	tokenInfo, err := s.firebase.VerifyToken(tokenStringArr[1])
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Check email already in use
	isExisted, err := s.repo.CheckExistedEmail(&tokenInfo.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if isExisted {
		return nil, http.StatusOK, errors.New(common.MessageErrorExistedEmail)
	}

	user := &entity.User{
		ID:             tokenInfo.UserID,
		Email:          tokenInfo.Email,
		Name:           tokenInfo.Name,
		OrganizationID: primitive.NilObjectID,
		Status:         "Active",
	}
	user.SetTime()

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return createdUser, http.StatusOK, nil
}

func (s *Service) GetUserByID(userID *string) (*entity.User, int, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (s *Service) GetUserByEmail(email *string) (*entity.User, int, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

func (s *Service) TokenToUser(token *string) (*entity.User, int, error) {
	tokenStringArr := strings.Split(*token, " ")
	tokenInfo, err := s.firebase.VerifyToken(tokenStringArr[1])
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user, err := s.repo.GetUserByEmail(&tokenInfo.Email)
	if err == mongo.ErrNoDocuments {
		newUser := &entity.User{
			ID:             tokenInfo.UserID,
			Email:          tokenInfo.Email,
			Name:           tokenInfo.Name,
			OrganizationID: primitive.NilObjectID,
			Status:         "Active",
		}
		newUser.SetTime()
		user, err = s.repo.CreateUser(newUser)
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

func (s *Service) UpsertUserFromFireBase(userID *string) (*entity.User, int, error) {
	fbUser, err := s.firebase.GetFirebaseUser(*userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	newUser := entity.User{}.NewUserByFsUser(fbUser)
	newUser.SetTime()

	if &fbUser.Organization != nil {
		org, err := s.orgRepo.GetOrgByTagName(&fbUser.Organization)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		if org != nil {
			newUser.OrganizationID = org.ID
		}
	}

	err = s.repo.UpsertUser(&newUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &newUser, http.StatusOK, nil
}

func (s *Service) UpdateRole(req *request.UpdateRoleRequest) (bool, int, error) {
	ok, err := s.repo.UpdateRole(&req.NewRole, &req.UserID)
	if err != nil || !ok {
		return false, http.StatusInternalServerError, err
	}

	user, err := s.repo.GetUserByID(&req.UserID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	err = s.firebase.UpdateUserDetails(user)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	return true, http.StatusOK, err
}

func (s *Service) UpdateOrgID(req *request.UpdateOrgRequest) (bool, int, error) {
	ok, err := s.repo.UpdateOrgID(&req.NewOrgID, &req.UserID)
	if err != nil || !ok {
		return false, http.StatusInternalServerError, err
	}

	user, err := s.repo.GetUserByID(&req.UserID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	user.Organization = req.NewOrgTagName

	err = s.firebase.UpdateUserDetails(user)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}

	return true, http.StatusOK, err
}

func (s *Service) UpdateUserDetails(userID *string, req *request.UpdateUserDetailsRequest) (bool, int, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	if user == nil {
		// Get user from Firebase
		fbUser, err := s.firebase.GetFirebaseUser(*userID)
		if err != nil {
			return false, http.StatusInternalServerError, err
		}
		newUser := entity.User{}.NewUserByFsUser(fbUser)
		newUser.SetTime()
		helper.MergeStructsField(&newUser, *req)

		if &fbUser.Organization != nil {
			org, err := s.orgRepo.GetOrgByTagName(&fbUser.Organization)
			if err != nil {
				return false, http.StatusInternalServerError, err
			}

			if org != nil {
				newUser.OrganizationID = org.ID
			}
		}

		err = s.repo.UpsertUser(&newUser)
		if err != nil {
			return false, http.StatusInternalServerError, err
		}
		return true, http.StatusOK, nil
	} else {
		ok, err := s.repo.UpdateUserDetails(userID, req)
		if err != nil {
			return false, http.StatusInternalServerError, err
		}
		return ok, http.StatusOK, nil
	}
}

func (s *Service) GetUserWithNoWallet() (*[]entity.User, int, error) {
	usersNoWallet, err := s.repo.GetUserWithNoWallet()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return usersNoWallet, http.StatusOK, nil
}
