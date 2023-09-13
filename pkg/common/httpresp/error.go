package httpresp

import "errors"

type ErrorKey string

var (
	ErrKeySystemInternalServer                 = errors.New("error.system.internal")
	ErrKeyAuthenticationNoPermission           = errors.New("error.authentication.no_permission")
	ErrKeyAuthenticationInvalidAuthTokenFormat = errors.New("error.authentication.invalid_auth_token_format")
	ErrKeyAuthenticationInvalidSignature       = errors.New("error.authentication.invalid_signature")
	ErrKeyHTTPValidatorsMissingRequiredField   = errors.New("error.http_validator.missing_required_field")
	ErrKeyHTTPValidatorsInvalidFieldType       = errors.New("error.http_validator.invalid_filed_type")
	ErrKeyHTTPValidatorsDecodeFail             = errors.New("error.http_validator.decode_fail")
	ErrKeyDatabaseNotFound                     = errors.New("error.database.not_found")
)

func NewError(key string) error {
	return errors.New(key)
}
