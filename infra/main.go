package main

import (
	"fmt"
	"os"
	"strings"
	"swiftbeaver-infra/config"
	"swiftbeaver-infra/stack"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiconfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func NewConfig(ctx *pulumi.Context) (*config.Config, error) {
	stackOrg := pulumiconfig.Require(ctx, fmt.Sprintf("%s:stackOrg", ctx.Project()))

	gcpProjectId := pulumiconfig.Require(ctx, "gcp:project")
	gcpRegion := pulumiconfig.Require(ctx, "gcp:region")

	cfg := &config.Config{
		ProjectName: ctx.Project(),
		StackOrg:    stackOrg,

		GCPProjectID: gcpProjectId,
		GCPRegion:    gcpRegion,
	}

	stackName := ctx.Stack()

	if stackName == "infra" {
		cfg.StackName = stackName
		cfg.Environment = config.GlobalEnvironment
	} else {
		isDev := strings.HasSuffix(stackName, "-dev")
		isProd := strings.HasSuffix(stackName, "-prod")

		if !isDev && !isProd {
			return nil, fmt.Errorf("unknown stack: %s", stackName)
		}

		stackNameParts := strings.Split(stackName, "-")

		cfg.StackName = strings.Join(stackNameParts[:len(stackNameParts)-1], "-")
		cfg.Environment = config.Environment(stackNameParts[len(stackNameParts)-1])

		if cfg.StackName == "infra" {
			cfg.InfraEnvConfig = &config.InfraEnvConfig{
				APIServiceURL: os.Getenv("API_SERVICE_URL"),
			}
		}
	}

	if !cfg.Valid() {
		return nil, fmt.Errorf("invalid config")
	}

	return cfg, nil
}

func RunStack(ctx *pulumi.Context, stack stack.Stack, err error) error {
	if err != nil {
		return err
	}

	if err := stack.Run(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg, err := NewConfig(ctx)
		if err != nil {
			return err
		}

		switch cfg.StackName {
		case "infra":
			if cfg.Environment == config.GlobalEnvironment {
				stack, err := stack.NewInfraStack(cfg)
				if err := RunStack(ctx, stack, err); err != nil {
					return err
				}
			} else {
				stack, err := stack.NewInfraEnvStack(cfg)
				if err := RunStack(ctx, stack, err); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unknown stack: %s", cfg.StackName)
		}

		return nil
	})
}
