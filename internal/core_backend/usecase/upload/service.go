package upload

import (
	"fmt"
	"mime/multipart"
	"net/http"

	config "backend-service/config/core_backend"
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/common/logger"
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/infrastructure/storage"

	"strconv"
	"sync"
	"time"
)

// Service struct
type Service struct {
	googleStorage *storage.GCPClient
}

// NewService create service
func NewService(gs *storage.GCPClient) *Service {
	return &Service{
		googleStorage: gs,
	}
}

func (s *Service) UploadImagesToGCPStorage(request *request.UploadImagesRequest) (*[]entity.Upload, int, error) {
	var wg sync.WaitGroup
	total := len(request.Images)

	var returnUpload = make([]entity.Upload, total)
	uploadResult := make(chan entity.Upload, total)
	for index, image := range request.Images {
		wg.Add(1)
		go func(img *multipart.FileHeader, idx int) {
			defer wg.Done()

			file, err := img.Open()
			if err != nil {
				uploadResult <- entity.Upload{
					ID:           idx,
					ImgPath:      "",
					MessageError: err.Error(),
					Error:        err,
				}
				return
			}

			mimeType := img.Header.Get("Content-Type")
			file.Seek(0, 0)

			fileName := img.Filename
			var (
				pendingGCP int
				gcpPath    string
			)
			for pendingGCP < 3 {
				// TODO: "user_id" will be replaced with real user ID
				gcpPath, err = s.googleStorage.UploadFile(
					file,
					fmt.Sprintf("%s/%s", config.C.GCP.StorageImagePath, "user_id"),
					fmt.Sprintf("%s-%s", strconv.FormatInt(time.Now().Unix(), 10), fileName),
					mimeType,
				)
				if err != nil {
					pendingGCP++
					if pendingGCP < 3 {
						time.Sleep(5 * time.Second)
						continue
					}
					logger.LogInfo(fmt.Sprintf("Fail to upload image %s to GCP Storage ", fileName))
				}
				pendingGCP = 3
			}

			if err != nil {
				uploadResult <- entity.Upload{
					ID:           idx,
					ImgPath:      "",
					MessageError: fmt.Sprintf("Upload image %s to GCP Storage fail: %v", fileName, err),
					Error:        err,
				}
				return
			}

			uploadResult <- entity.Upload{
				ID:           idx,
				ImgPath:      gcpPath,
				MessageError: "",
				Error:        nil,
			}
		}(image, index)
	}
	wg.Wait()

	var suc, fail int
	for result := range uploadResult {
		if result.Error != nil {
			logger.LogError(fmt.Sprintf("Got error %v when upload images, message: %s", result.Error, result.MessageError))
			fail++
			continue
		}
		suc++

		returnUpload[result.ID] = result
		if suc+fail == total {
			close(uploadResult)
		}
	}
	logger.LogSuccess(fmt.Sprintf("Upload images to GCP Storage: %d success and %d failure!", suc, fail))

	return &returnUpload, http.StatusOK, nil
}
