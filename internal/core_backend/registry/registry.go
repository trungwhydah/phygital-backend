// Package registry Common registration
package registry

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/middleware"
	"backend-service/internal/core_backend/infrastructure/callers"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"backend-service/internal/core_backend/infrastructure/storage"
	validation "backend-service/internal/core_backend/infrastructure/validator"
	"backend-service/internal/core_backend/usecase/nft"
)

type interactor struct {
	mongo     *mongo.Database
	validator *validator.Validate
	caller    *callers.Caller
	firebase  *firebase.FirebaseClient
	gStorage  *storage.GCPClient
}

// Interactor Interactor interface
type Interactor interface {
	NewAppHandler() handler.AppHandler
	NewMiddlewareServices() middleware.MidddlewareServices
	NewNFTGlobalService() *nft.Service
}

// NewInteractor Constructs new interactor
func NewInteractor(mongo *mongo.Database, v *validator.Validate, c *callers.Caller, fb *firebase.FirebaseClient, gs *storage.GCPClient) Interactor {
	return &interactor{mongo: mongo, validator: v, caller: c, firebase: fb, gStorage: gs}
}

// NewAppHandler register all app handler
func (i *interactor) NewAppHandler() handler.AppHandler {
	return handler.AppHandler{
		DigitalAssetHandler: i.NewDigitalAssetHandler(),
		DummyHandler:        i.NewDummyHandler(),
		UserHandler:         i.NewUserHandler(),
		ScanHandler:         i.NewScanHandler(),
		MappingHandler:      i.NewMappingHandler(),
		ProductItemHandler:  i.NewProductItemHandler(),
		TagHandler:          i.NewTagHandler(),
		WebPageHandler:      i.NewWebPageHandler(),
		ProductHandler:      i.NewProductHandler(),
		OrganizationHandler: i.NewOrganizationHandler(),
		SessionHandler:      i.NewSessionHandler(),
		UploadHandler:       i.NewUploadHandler(),
		TemplateHandler:     i.NewTemplateHandler(),
		PubsubHandler:       i.NewPubsubHandler(),
		AuthorHandler:       i.NewAuthorHandler(),
	}
}

// NewCustomValidator register custom validator
func (i *interactor) NewCustomValidator() validation.CustomValidator {
	return validation.NewCustomValidator(i.validator)
}

func (i *interactor) NewCaller() *callers.Caller {
	return callers.NewCaller()
}

func (i *interactor) NewMiddlewareServices() middleware.MidddlewareServices {
	return middleware.NewMiddlewareServices(i.firebase, i.NewProductItemRepository(), i.NewProductRepository())
}

func (i *interactor) NewNFTGlobalService() *nft.Service {
	return i.NewNFTService()
}
