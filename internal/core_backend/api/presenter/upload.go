package presenter

import (
	"backend-service/internal/core_backend/entity"
)

// UploadResponse data struct
type UploadResponse struct {
	ImagePaths []string `json:"image_paths"`
}

// presenterUpload struct
type PresenterUpload struct{}

// presenterUpload interface
type ConvertUpload interface {
	ResponseUpload(*[]entity.Upload) *UploadResponse
}

// NewPresenterUpload Constructs presenter
func NewPresenterUpload() ConvertUpload {
	return &PresenterUpload{}
}

// Return property data response
func (pp *PresenterUpload) ResponseUpload(resultUpload *[]entity.Upload) *UploadResponse {
	var response UploadResponse
	for i := 0; i < len(*resultUpload); i++ {
		response.ImagePaths = append(response.ImagePaths, (*resultUpload)[i].ImgPath)
	}

	return &response
}
