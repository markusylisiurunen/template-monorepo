package stack

import (
	"swiftbeaver-infra/config"
	"swiftbeaver-infra/recipe"
	"swiftbeaver-infra/util"

	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InfraStack struct {
	cfg     *config.Config
	affixer *util.Affixer
}

func (stack *InfraStack) setupRegistry(ctx *pulumi.Context) error {
	registry, err := container.NewRegistry(ctx, "registry",
		&container.RegistryArgs{
			Location: pulumi.String("EU"),
		},
	)
	if err != nil {
		return err
	}

	registryURL := registry.ID().ToStringOutput().ApplyT(
		func(_ string) (string, error) {
			repository, err := container.GetRegistryRepository(ctx, nil)
			if err != nil {
				return "", err
			}

			return repository.RepositoryUrl, nil
		},
	)

	ctx.Export(stack.affixer.PrefixOutput(config.RegistryURLOutput), registryURL)

	return nil
}

func (stack *InfraStack) setupDatabase(ctx *pulumi.Context) error {
	instance, err := recipe.NewDatabaseInstance(ctx, "primary", &recipe.DatabaseInstanceArgs{
		Name:               pulumi.String(stack.affixer.PrefixResource("primary")),
		Tier:               pulumi.String("db-f1-micro"),
		DiskSize:           pulumi.Int(10),
		RootPasswordSecret: pulumi.String(stack.affixer.PrefixResource("primary-root-password")),
	})
	if err != nil {
		return err
	}

	ctx.Export(stack.affixer.PrefixOutput(config.DatabaseInstanceIDOutput), instance.Instance.ID())
	ctx.Export(stack.affixer.PrefixOutput(config.DatabaseRootUserNameOutput), instance.RootUser.Name)
	ctx.Export(stack.affixer.PrefixOutput(config.DatabaseRootPasswordSecretOutput), instance.RootPasswordSecret.SecretId)

	return nil
}

func (stack *InfraStack) Run(ctx *pulumi.Context) error {
	if err := stack.setupRegistry(ctx); err != nil {
		return err
	}

	if err := stack.setupDatabase(ctx); err != nil {
		return err
	}

	return nil
}

func NewInfraStack(cfg *config.Config) (Stack, error) {
	return &InfraStack{
		cfg:     cfg,
		affixer: util.NewAffixer(cfg),
	}, nil
}
