package app

import (
	cmconfig "backend-service/config/common"
	config "backend-service/config/marketplace"
	cmdomain "backend-service/internal/common/domain"
	cmrepo "backend-service/internal/common/repo"
	"backend-service/internal/marketplace/api"
	"backend-service/internal/marketplace/api/handler"
	"backend-service/internal/marketplace/api/middleware"
	"backend-service/internal/marketplace/domain"
	"backend-service/internal/marketplace/repo"
	"backend-service/pkg/common/logger"
	"backend-service/pkg/common/mongo"
	"backend-service/pkg/common/msgtranslate"
	"backend-service/pkg/marketplace/firebase"
	"backend-service/pkg/marketplace/firebase/storage"

	"go.uber.org/fx"
)

var InternalOptions = fx.Options(
	// Common Config
	fx.Provide(cmconfig.NewConfig),

	// Config
	fx.Provide(config.NewConfig),

	// Server
	fx.Provide(NewServer),

	// Router
	fx.Provide(api.NewRouter),

	// Controller
	handler.Module,

	// Middleware
	middleware.Module,

	// Use Case
	domain.Module,

	// Repo
	repo.Module,

	// Common Repo
	cmrepo.Module,

	// Common Domain
	cmdomain.Module,
)

var PackageOptions = fx.Options(
	// Mongo
	fx.Provide(mongo.New),

	// Firebase
	fx.Provide(firebase.NewApps),

	// Storage
	fx.Provide(storage.NewBucketHandler),

	// Logger
	fx.Provide(logger.Init),

	// Msg translate
	fx.Invoke(msgtranslate.Init),
)
