package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	config "backend-service/config/core_backend"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

type GCPClient struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

func NewGCPClient() *GCPClient {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile("./service-account-file.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil
	}

	return &GCPClient{
		client:     client,
		bucketName: config.C.GCP.StorageBucketName,
		projectID:  config.C.GCP.StorageProjectID,
	}
}

// UploadFile uploads a file
func (c *GCPClient) UploadFile(file multipart.File, uploadPath, fileName, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	wc := c.client.Bucket(c.bucketName).Object(
		fmt.Sprintf("%s/%s",
			uploadPath,
			fileName)).
		NewWriter(ctx)
	wc.ContentType = contentType
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return fmt.Sprintf("%s/%s/%s/%s", config.C.GCP.StorageDomain, c.bucketName, uploadPath, fileName), nil
}
