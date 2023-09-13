package entity

type Upload struct {
	ID           int
	ImgPath      string
	MessageError string
	Error        error
}

type UploadResultList struct {
	UploadList []Upload `json:"upload_list"`
}
