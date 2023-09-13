package helper

import (
	"regexp"

	"backend-service/internal/core_backend/common/logger"
)

func ValidateOrganizationTagName(orgTagName string) (bool, error) {
	// Regular expression pattern to validate organization tag name
	pattern := "^[a-zA-Z0-9_-]+$"
	found, err := regexp.MatchString(pattern, orgTagName)
	if err != nil {
		logger.LogError(err.Error())
		return false, err
	}
	return found, err
}
