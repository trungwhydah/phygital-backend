package request

type CreateTagRequest struct {
	HardwareID     string `form:"hardware_id"`
	TagID          string `form:"tag_id" validate:"required"`
	TagType        string `form:"tag_type"`
	EncryptMode    string `form:"encrypt_mode"`
	RawData        string `form:"raw_data"`
	ScanCounter    int    `form:"scan_counter"`
	OrganizationID string `form:"org_id" validate:"required"`
}
