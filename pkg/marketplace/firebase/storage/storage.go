package storage

import (
	"context"
	"fmt"
	"io"

	config "backend-service/config/marketplace"

	"cloud.google.com/go/storage"
	fbstorage "firebase.google.com/go/storage"
)

func UploadFile(ctx context.Context, client *fbstorage.Client, path string, file io.Reader) (string, error) {
	handler, err := client.Bucket("raramuriapp-dev")
	if err != nil {
		return "", fmt.Errorf("client.Bucket(): %s", err)
	}

	o := handler.Object(path)
	objAppliedCond := o.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer.
	wc := objAppliedCond.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy(): %s", err)
	}

	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("wc.Close(): %s", err)
	}

	// publish for anonymous users
	acl := o.ACL()
	if err = acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("acl.Set(): %s", err)
	}

	// get media URL
	attr, err := o.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("o.Attrs(): %s", err)
	}

	return attr.MediaLink, nil
}

func NewBucketHandler(cfg *config.Config, client *fbstorage.Client) (*storage.BucketHandle, error) {
	return client.Bucket(cfg.FirebaseStorage.BucketName)
}
