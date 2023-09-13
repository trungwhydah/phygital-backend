package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App             `yaml:"app"`
		HTTP            `yaml:"http"`
		Log             `yaml:"logger"`
		Mongo           `yaml:"mongo"`
		FirebaseStorage `yaml:"firebase_storage"`
	}

	// App specific general information of service.
	//
	// Name of service.
	// Env type of service, can be dev or production, if value is production swagger will be ignored.
	// logger will be applied production mode.
	App struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Env     string `yaml:"env" env:"APP_ENV"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// HTTP specific information of http rest api.s
	//
	// Port default running at 8080.
	HTTP struct {
		Host string `yaml:"host" env:"HTTP_HOST"`
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	// Log specific level of logger using in service.
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
		Prod  bool   `yaml:"prod" env:"LOG_PROD"`
	}

	// Mongo specific mongo connection info.
	Mongo struct {
		DBName  string `yaml:"db_name"  env:"MONGO_DB_NAME"`
		ConnURI string `yaml:"conn_uri" env:"MONGO_CONN_URI"`
	}

	FirebaseStorage struct {
		BucketName string `yaml:"bucket_name" env:"FIREBASE_STORAGE_BUCKET"`
	}
)

const EnvProd = "production"

// NewConfig returns app config.
//
// read config from config.yml if env is dev.
// read config from config.production.yml and overwrite config.yml if env is production.
// overwrite yml config by variables in .env file.
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = cleanenv.ReadConfig(
		"./config/common/config.yml",
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"read common config error: %w",
			err,
		)
	}

	err = cleanenv.ReadConfig(
		"./config/marketplace/config.yml",
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"read base config error: %w",
			err,
		)
	}

	if cfg.App.Env == EnvProd {
		// overwrite some values from /config/config.production.yml
		err = cleanenv.ReadConfig(
			"./config/marketplace/config.production.yml",
			cfg,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"read production config error: %w",
				err,
			)
		}
	}

	// lastly, overwrite value from environment variable
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
