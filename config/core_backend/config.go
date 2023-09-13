package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"backend-service/internal/core_backend/common/logger"
)

type config struct {
	Server struct {
		SimultaneousConnection int `env:"SIMULTANEOUS_CONNECTION" envDefault:"50"`
		SessionTimeoutInSecond int `env:"SESSION_TIMEOUT_IN_SECOND" envDefault:"10"`
	}
	Mongo struct {
		DatabaseName     string `env:"MONGO_DATABASE_NAME" envDefault:""`
		IntervalTime     int    `env:"MONGO_INTERVAL_TIME" envDefault:"3"`
		MaxRetry         int    `env:"MONGO_MAX_RETRY" envDefault:"10"`
		MongoURLDbString string `env:"MONGO_URL_DB_STRING"`
	}
	Domains struct {
		AuthenticateNFCDomain string `env:"AUTHENTICATE_NFC_DOMAIN"`
		WebpageDomain         string `env:"WEBPAGE_DOMAIN"`
		ScanErrorPage         string `env:"SCAN_ERROR_PAGE"`
	}
	Firebase struct {
		FirebaseProjectID string `env:"FIREBASE_PROJECT_ID"`
	}
	GCP struct {
		StorageBucketName string `env:"STORAGE_BUCKET_NAME"`
		StorageProjectID  string `env:"STORAGE_PROJECT_ID"`
		StorageDomain     string `env:"GCP_STORAGE_DOMAIN"`
		StorageImagePath  string `env:"GCP_STORAGE_IMAGE_PATH"`
	}
	Pubsub struct {
		AuthenToken string `env:"PUBSUB_AUTHTOKEN_TOKEN"`
	}
	Cheat struct {
		ProductID string `env:"PRODUCT_ID"`
	}
	NFT struct {
		BASE_NET_URL                    string `env:"BASE_NET_URL"`
		PRIVATE_KEY                     string `env:"PRIVATE_KEY"`
		WATCH_TRANSACTION_MAX_RETRY     int    `env:"WATCH_TRANSACTION_MAX_RETRY"`
		WATCH_TRANSACTION_INTERVAL_TIME int    `env:"WATCH_TRANSACTION_INTERVAL_TIME"`
	}
	Wallet struct {
		WALLET_CLIENT_ID   string `env:"WALLET_CLIENT_ID"`
		WALLET_CLIENT_KEY  string `env:"WALLET_CLIENT_KEY"`
		WALLET_DOMAIN_V1   string `env:"WALLET_DOMAIN_V1"`
		WALLET_API_TIMEOUT int    `env:"WALLET_API_TIMEOUT"`
	}
}

// C config struct
var C config

// LoadConfig load config from environment and parse to struct
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError("Application config parsing failed: " + err.Error())
		logger.LogInfo("Could not run application!")
		return
	}

	cfg := &C
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		logger.LogError("Application config parsing failed: " + err.Error())
		logger.LogInfo("Could not run application!")
		return
	}

	logger.LogSuccess("Load Config Successfully!")
}
