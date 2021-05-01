package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
	"github.com/spf13/viper"
)

type Env int

const (
	Development Env = iota + 1
	Production
)

type Config struct {
	DatabaseURL string `validate:"required"`
	Env         Env    `validate:"required,gt=0"`
	ServerHost  string `validate:"required"`
	ServerPort  int    `validate:"required"`
}

func Load() error {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		logger.Default.Debugf("config file (.env) was not found")
	}

	viper.SetEnvPrefix("SWIFTBEAVER")
	viper.AutomaticEnv()

	portBindErr := viper.BindEnv("ServerPort", "PORT")
	if portBindErr != nil {
		logger.Default.Warnf("failed to bind 'ServerPort' to 'PORT'")
	}

	// Heroku configures the database connection as DATABASE_URL env variable
	databaseBindErr := viper.BindEnv("DatabaseURL", "DATABASE_URL")
	if databaseBindErr != nil {
		logger.Default.Warnf("failed to bind 'DatabaseURL' to 'DATABASE_URL'")
	}

	viper.SetDefault("ServerHost", "0.0.0.0")

	return nil
}

func Get() (*Config, error) {
	envMap := map[string]Env{
		"development": Development,
		"production":  Production,
	}

	config := &Config{
		DatabaseURL: viper.GetString("DatabaseURL"),
		Env:         envMap[viper.GetString("Env")],
		ServerHost:  viper.GetString("ServerHost"),
		ServerPort:  viper.GetInt("ServerPort"),
	}

	if err := validator.New().Struct(config); err != nil {
		for _, err = range err.(validator.ValidationErrors) {
			logger.Default.Errorf(err.Error())
		}

		return nil, err
	}

	return config, nil
}
