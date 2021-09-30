package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Environment string

const (
	GlobalEnvironment      Environment = "global"
	DevelopmentEnvironment Environment = "dev"
	ProductionEnvironment  Environment = "prod"
)

type InfraEnvConfig struct {
	APIServiceURL string `validate:"required"`
}

type Config struct {
	Environment Environment `validate:"required,oneof=global dev prod"`
	ProjectName string      `validate:"required"`
	StackName   string      `validate:"required"`
	StackOrg    string      `validate:"required"`

	GCPProjectID string `validate:"required"`
	GCPRegion    string `validate:"required"`

	InfraEnvConfig *InfraEnvConfig
}

func (cfg *Config) Valid() bool {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		for _, err = range err.(validator.ValidationErrors) {
			fmt.Println(err.Error())
		}

		return false
	}

	if cfg.StackName == "infra" && cfg.Environment != GlobalEnvironment {
		if cfg.InfraEnvConfig == nil {
			fmt.Println("missing environment-specific infra config")
			return false
		}

		if err := validate.Struct(cfg); err != nil {
			for _, err = range err.(validator.ValidationErrors) {
				fmt.Println(err.Error())
			}

			return false
		}
	}

	return true
}
