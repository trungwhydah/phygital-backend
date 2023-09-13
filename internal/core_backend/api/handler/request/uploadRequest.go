package request

import "mime/multipart"

type UploadImagesRequest struct {
	Images []*multipart.FileHeader `form:"-"`
}
