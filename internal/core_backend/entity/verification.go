package entity

type Verification struct {
	BaseModel `bson:"inline"`
	TagID     string `bson:"tag_id"`
	IsValid   bool   `bson:"is_valid"`
	Nonce     int    `bson:"nonce"`
}

// CollectionName Collection name of Verification
func (Verification) CollectionName() string {
	return "verifications"
}

type Scan struct {
	UID         string `bson:"uid" json:"uid"`
	TagID       string `bson:"tag_id" json:"tag_id"`
	ScanCounter int    `bson:"scan_counter" json:"scan_counter"`
	EncMode     string `bson:"enc_mode" json:"enc_mode"`
	Error       error  `bson:"error" json:"error"`
	StatusCode  int    `bson:"status_code" json:"status_code"`
}
