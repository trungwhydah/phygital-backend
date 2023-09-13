package firebase

import (
	"context"

	"backend-service/pkg/common/logger"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"go.uber.org/fx"
)

// Apps specific all client of firebase service used in service.
type Apps struct {
	fx.Out
	Auth    *auth.Client
	Storage *storage.Client
}

// NewApps returns Apps, which contains all instances firebase used in service.
func NewApps() Apps {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		logger.Fatal("Fail to init Firebase client", "error", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal("Fail to init Firebase Auth client", "error", err)
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		logger.Fatal("Fail to init Storage client", "error", err)
	}

	return Apps{Auth: authClient, Storage: storageClient}
}
