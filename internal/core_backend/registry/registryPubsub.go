package registry

import (
	"backend-service/internal/core_backend/api/handler"
)

// NewPubsubHandler
func (i *interactor) NewPubsubHandler() handler.PubsubHandler {
	return handler.NewPubsubHandler(i.NewUserService(), i.NewCustomValidator())
}
