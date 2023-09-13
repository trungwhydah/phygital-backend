package registry

import (
	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/usecase/upload"
)

// NewUploadService new upload service
func (i *interactor) NewUploadService() *upload.Service {
	return upload.NewService(i.gStorage)
}

// NewUploadPresenter
func (i *interactor) NewUploadPresenter() presenter.ConvertUpload {
	return presenter.NewPresenterUpload()
}

// NewUploadHandler
func (i *interactor) NewUploadHandler() handler.UploadHandler {
	return handler.NewUploadHandler(i.NewUploadService(), i.NewUploadPresenter(), i.NewCustomValidator())
}
