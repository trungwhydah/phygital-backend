package scan

import (
	"encoding/json"
	"net/http"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/infrastructure/callers"
)

// Service struct
type Service struct {
	caller *callers.Caller
	repo   Repository
}

// NewService create service
func NewService(c *callers.Caller, r Repository) *Service {
	return &Service{
		caller: c,
		repo:   r,
	}
}

func (s *Service) ProcessScan(request *request.ScanRequest) (*entity.Scan, int, error) {
	query := make(map[string]string)
	query["picc_data"] = request.PiccData
	query["enc"] = request.Enc
	query["cmac"] = request.Cmac

	responseBody, err := s.caller.CallGetMethod(config.C.Domains.AuthenticateNFCDomain, "api/tag", query)
	if err != nil {
		logger.LogError("[DEBUG] - 3 - Got error while calling to Authentication" + err.Error())
		return nil, http.StatusInternalServerError, err
	}

	var result entity.Scan
	err = json.NewDecoder(responseBody).Decode(&result)
	if err != nil {
		logger.LogError("[DEBUG] - 4 - Fail to parse response body" + err.Error())
		return nil, http.StatusInternalServerError, err
	}
	responseBody.Close()

	if result.Error != nil {
		logger.LogError("[DEBUG] - 5 - Got error while verifying tag: " + result.Error.Error())
		return &result, result.StatusCode, result.Error
	}

	return &result, http.StatusOK, nil
}
