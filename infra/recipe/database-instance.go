package recipe

import (
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/sql"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DatabaseInstance struct {
	pulumi.ResourceState

	Instance           *sql.DatabaseInstance
	RootUser           *sql.User
	RootPasswordSecret *secretmanager.Secret
}

type DatabaseInstanceArgs struct {
	Name               pulumi.StringInput
	Tier               pulumi.StringInput
	DiskSize           pulumi.IntInput
	RootPasswordSecret pulumi.StringInput
}

func NewDatabaseInstance(
	ctx *pulumi.Context,
	name string,
	args *DatabaseInstanceArgs,
	opts ...pulumi.ResourceOption,
) (*DatabaseInstance, error) {
	instance := &DatabaseInstance{}

	if err := ctx.RegisterComponentResource("DatabaseInstanceRecipe", name, instance, opts...); err != nil {
		return nil, err
	}

	// create the Cloud SQL instance
	hash, err := random.NewRandomId(ctx, "instance-hash",
		&random.RandomIdArgs{
			ByteLength: pulumi.Int(4),
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	physicalInstance, err := sql.NewDatabaseInstance(ctx, "instance",
		&sql.DatabaseInstanceArgs{
			DatabaseVersion: pulumi.String("POSTGRES_13"),
			Name:            pulumi.Sprintf("%s-%s", args.Name, hash.Hex),

			Settings: &sql.DatabaseInstanceSettingsArgs{
				AvailabilityType: pulumi.String("ZONAL"),
				DiskAutoresize:   pulumi.Bool(true),
				DiskSize:         args.DiskSize,
				Tier:             args.Tier,

				BackupConfiguration: &sql.DatabaseInstanceSettingsBackupConfigurationArgs{
					Enabled:                    pulumi.Bool(true),
					PointInTimeRecoveryEnabled: pulumi.Bool(true),
				},

				InsightsConfig: &sql.DatabaseInstanceSettingsInsightsConfigArgs{
					QueryInsightsEnabled: pulumi.Bool(true),
				},
			},
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	instance.Instance = physicalInstance

	// create the credentials for the root user
	rootPassword, err := random.NewRandomPassword(ctx, "root-password",
		&random.RandomPasswordArgs{
			Length:          pulumi.Int(32),
			OverrideSpecial: pulumi.String("_%@"),
			Special:         pulumi.Bool(true),
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	rootPasswordSecret, err := secretmanager.NewSecret(ctx, "root-password",
		&secretmanager.SecretArgs{
			SecretId: args.RootPasswordSecret,

			Replication: &secretmanager.SecretReplicationArgs{
				Automatic: pulumi.Bool(true),
			},
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	_, err = secretmanager.NewSecretVersion(ctx, "root-password",
		&secretmanager.SecretVersionArgs{
			Enabled:    pulumi.Bool(true),
			Secret:     rootPasswordSecret.ID(),
			SecretData: rootPassword.Result,
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	instance.RootPasswordSecret = rootPasswordSecret

	// create the actual root user
	rootUser, err := sql.NewUser(ctx, "root",
		&sql.UserArgs{
			Instance: physicalInstance.Name,
			Name:     pulumi.String("root"),
			Password: rootPassword.Result,
		},
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, err
	}

	instance.RootUser = rootUser

	ctx.RegisterResourceOutputs(instance, pulumi.Map{})

	return instance, nil
}
