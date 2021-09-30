package recipe

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/sql"
	"github.com/pulumi/pulumi-postgresql/sdk/v3/go/postgresql"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Database struct {
	pulumi.ResourceState

	Database       *sql.Database
	PasswordSecret *secretmanager.Secret
}

type DatabaseArgs struct {
	RootPasswordSecretName string

	InstanceID         pulumi.IDInput
	DatabaseName       pulumi.StringInput
	PasswordSecretName pulumi.StringInput
}

func NewDatabase(
	ctx *pulumi.Context,
	name string,
	args *DatabaseArgs,
	opts ...pulumi.ResourceOption,
) (*Database, error) {
	database := &Database{}

	if err := ctx.RegisterComponentResource("DatabaseRecipe", name, database, opts...); err != nil {
		return nil, err
	}

	// create the credentials for the database user
	password, err := random.NewRandomPassword(ctx, fmt.Sprintf("%s-password", name),
		&random.RandomPasswordArgs{
			Length:          pulumi.Int(32),
			OverrideSpecial: pulumi.String("_%@"),
			Special:         pulumi.Bool(true),
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	passwordSecret, err := secretmanager.NewSecret(ctx, fmt.Sprintf("%s-password", name),
		&secretmanager.SecretArgs{
			SecretId: args.PasswordSecretName,

			Replication: &secretmanager.SecretReplicationArgs{
				Automatic: pulumi.Bool(true),
			},
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	database.PasswordSecret = passwordSecret

	_, err = secretmanager.NewSecretVersion(ctx, fmt.Sprintf("%s-password", name),
		&secretmanager.SecretVersionArgs{
			Enabled:    pulumi.Bool(true),
			Secret:     passwordSecret.ID(),
			SecretData: password.Result,
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	// create the actual database in the instance
	instance, err := sql.GetDatabaseInstance(
		ctx,
		fmt.Sprintf("instance-for-%s", name),
		args.InstanceID,
		nil,
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	pgDatabase, err := sql.NewDatabase(ctx, name,
		&sql.DatabaseArgs{
			Instance: instance.Name,
			Name:     args.DatabaseName,
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	database.Database = pgDatabase

	// create the user in Postgres
	rootPasswordVersion := "latest"

	rootPassword, err := secretmanager.LookupSecretVersion(ctx,
		&secretmanager.LookupSecretVersionArgs{
			// FIXME: why does this require the value to be a regular string?
			Secret:  args.RootPasswordSecretName,
			Version: &rootPasswordVersion,
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	pg, err := postgresql.NewProvider(ctx, name,
		&postgresql.ProviderArgs{
			Host:      instance.PublicIpAddress,
			Port:      pulumi.Int(5432),
			Username:  pulumi.String("root"),
			Password:  pulumi.String(rootPassword.SecretData),
			Database:  pulumi.String("postgres"),
			Superuser: pulumi.Bool(false),
		},
		pulumi.Parent(database),
	)
	if err != nil {
		return nil, err
	}

	// disable all access to the created database from the `public` schema
	revokePublic, err := postgresql.NewGrant(ctx, fmt.Sprintf("%s-revoke-public-privileges", name),
		&postgresql.GrantArgs{
			Role:       pulumi.String("public"),
			Database:   pgDatabase.Name,
			ObjectType: pulumi.String("database"),
			Privileges: pulumi.StringArray{},
		},
		pulumi.Parent(database),
		pulumi.Provider(pg),
	)
	if err != nil {
		return nil, err
	}

	// create the new role (user) and grant access to only the created database
	createRole, err := postgresql.NewRole(ctx, fmt.Sprintf("%s-role", name),
		&postgresql.RoleArgs{
			Name:     args.DatabaseName,
			Password: password.Result,
			Login:    pulumi.Bool(true),
		},
		pulumi.Parent(database),
		pulumi.Provider(pg),
		pulumi.DependsOn([]pulumi.Resource{revokePublic}),
	)
	if err != nil {
		return nil, err
	}

	_, err = postgresql.NewGrant(ctx, fmt.Sprintf("%s-privileges", name),
		&postgresql.GrantArgs{
			Role:       args.DatabaseName,
			Database:   pgDatabase.Name,
			ObjectType: pulumi.String("database"),

			Privileges: pulumi.StringArray{
				pulumi.String("ALL"),
			},
		},
		pulumi.Parent(database),
		pulumi.Provider(pg),
		pulumi.DependsOn([]pulumi.Resource{createRole}),
	)
	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(database, pulumi.Map{})

	return database, nil
}
