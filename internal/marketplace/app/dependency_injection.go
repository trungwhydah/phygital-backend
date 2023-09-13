package app

import (
	cmconfig "backend-service/config/common"
	config "backend-service/config/marketplace"
	cmdomain "backend-service/internal/common/domain"
	cmrepo "backend-service/internal/common/repo"
	"backend-service/internal/marketplace/api/restful"
	"backend-service/internal/marketplace/api/restful/security"
	v1 "backend-service/internal/marketplace/api/restful/v1"
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
	fx.Provide(v1.NewRouter),

	// Controller
	restful.Module,

	// Security
	security.Module,

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
