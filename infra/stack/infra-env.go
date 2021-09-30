package stack

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"swiftbeaver-infra/config"
	"swiftbeaver-infra/recipe"
	"swiftbeaver-infra/util"
	"text/template"

	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InfraEnvStack struct {
	cfg     *config.Config
	affixer *util.Affixer
}

func (stack *InfraEnvStack) readOpenAPIConfig() ([]byte, error) {
	content, err := os.ReadFile(path.Join("static", "openapi.yaml"))
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"APIServiceURL": stack.cfg.InfraEnvConfig.APIServiceURL,
	}

	var buf bytes.Buffer

	if err := template.Must(template.New("OpenAPI").Parse(string(content))).Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (stack *InfraEnvStack) setupAPIGateway(ctx *pulumi.Context) error {
	api, err := apigateway.NewApi(ctx, "api",
		&apigateway.ApiArgs{
			ApiId: pulumi.String(stack.affixer.AffixResource("api")),
		},
	)
	if err != nil {
		return err
	}

	openAPIConfig, err := stack.readOpenAPIConfig()
	if err != nil {
		return err
	}

	openAPIConfigHash := sha1.New()
	if _, err := openAPIConfigHash.Write(openAPIConfig); err != nil {
		return err
	}

	apiConfig, err := apigateway.NewApiConfig(ctx, "api",
		&apigateway.ApiConfigArgs{
			Api: api.ApiId,

			ApiConfigId: pulumi.String(
				fmt.Sprintf("config-%s", hex.EncodeToString(openAPIConfigHash.Sum(nil))),
			),

			OpenapiDocuments: apigateway.ApiConfigOpenapiDocumentArray{
				apigateway.ApiConfigOpenapiDocumentArgs{
					Document: apigateway.ApiConfigOpenapiDocumentDocumentArgs{
						Path:     pulumi.String("spec.yaml"),
						Contents: pulumi.String(base64.StdEncoding.EncodeToString(openAPIConfig)),
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	gateway, err := apigateway.NewGateway(ctx, "api",
		&apigateway.GatewayArgs{
			ApiConfig: apiConfig.ID(),
			GatewayId: pulumi.String(stack.affixer.AffixResource("api")),
			Region:    pulumi.String("europe-west1"),
		},
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "15m"}),
	)
	if err != nil {
		return err
	}

	ctx.Export(stack.affixer.PrefixOutput(config.APIGatewayURLOutput), gateway.DefaultHostname)

	return nil
}

func (stack *InfraEnvStack) setupDatabases(ctx *pulumi.Context) error {
	infraStackName := fmt.Sprintf("%s/%s/infra", stack.cfg.StackOrg, stack.cfg.ProjectName)

	infraStack, err := pulumi.NewStackReference(ctx, infraStackName, nil)
	if err != nil {
		return err
	}

	instanceIdOutput := stack.affixer.PrefixOutput(config.DatabaseInstanceIDOutput)
	instanceId := infraStack.GetIDOutput(pulumi.String(instanceIdOutput))

	// create the databases
	names := []string{"api", "worker"}

	for _, name := range names {
		db, err := recipe.NewDatabase(ctx, name, &recipe.DatabaseArgs{
			RootPasswordSecretName: stack.affixer.PrefixResource("primary-root-password"),

			InstanceID:         instanceId,
			DatabaseName:       pulumi.String(name),
			PasswordSecretName: pulumi.String(stack.affixer.AffixResource(fmt.Sprintf("%s-database-password", name))),
		})
		if err != nil {
			return err
		}

		outputName := fmt.Sprintf("%s::%s", config.DatabasePasswordSecretOutput, name)

		ctx.Export(stack.affixer.PrefixOutput(outputName), db.PasswordSecret.SecretId)
	}

	return nil
}

func (stack *InfraEnvStack) Run(ctx *pulumi.Context) error {
	if err := stack.setupAPIGateway(ctx); err != nil {
		return err
	}

	if err := stack.setupDatabases(ctx); err != nil {
		return err
	}

	return nil
}

func NewInfraEnvStack(cfg *config.Config) (Stack, error) {
	return &InfraEnvStack{
		cfg:     cfg,
		affixer: util.NewAffixer(cfg),
	}, nil
}
