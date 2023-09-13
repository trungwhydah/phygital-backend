package request

type ScanRequest struct {
	PiccData string `form:"picc_data" validate:"required"`
	Enc      string `form:"enc"`
	Cmac     string `form:"cmac"`
}

type TapRequest struct {
	TagID string `validate:"required" swaggerignore:"true"`
}
