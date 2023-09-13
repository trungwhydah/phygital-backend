package request

import "strconv"

type GetAssetByCollectionRequest struct {
	CollectionID string
}

type GetDigitalMetadataWithIDRequest struct {
	OrgTagName string `validate:"required"`
	TokenID    string `validate:"required,numeric"`
}

func (r *GetDigitalMetadataWithIDRequest) TokenIDToInt() int {
	tokenIDInInt, _ := strconv.Atoi(r.TokenID)

	return tokenIDInInt
}
