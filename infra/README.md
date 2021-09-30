# Infrastructure

1. [Quickstart](#quickstart)
2. [Stacks](#stacks)
   1. [infra](#infra)
   2. [infra-env](#infra-env)

## Quickstart

To get started, go to Google Cloud Platform's (GCP) console and create a new empty project. Once
created, run the following on your machine where `$PROJECT_ID` is the GCP project ID of the new
project you just created.

```bash
gcloud config set project $PROJECT_ID
gcloud auth application-default login
```

Next, edit the `Pulumi.{infra,infra-dev}.yaml` files to have the correct information. You need to
update the values for at least `gcp:project` (GCP project ID) and `swiftbeaver-infra:stackOrg`
(Pulumi's organisation).

You are now ready to enable the necessary GCP APIs. To do so, run the following.

```bash
services=(
  apigateway.googleapis.com
  compute.googleapis.com
  containerregistry.googleapis.com
  secretmanager.googleapis.com
  servicecontrol.googleapis.com
  sql-component.googleapis.com
  sqladmin.googleapis.com
)

for service in "${services[@]}"; do
  gcloud services enable $service
done
```

> **Note:** It may take a few minutes for the APIs to be enabled. Be patient.

Finally, you can start by deploying the bottom-most `infra` stack to GCP.

```bash
pulumi stack init infra
pulumi stack select infra

pulumi up
```

The next layer is the environment-specific infra stack. Before deploying the stack, make sure your
current IP address is allowed to connect to the Cloud SQL instance created by the `infra` stack.
Otherwise, the PostgreSQL provider is not able to connect, create the roles, and grant the necessary
privileges.

You can allow your IP either from the GCP console manually or by running the following.

```bash
SQL_INSTANCE=$(pulumi stack output swiftbeaver-infra::databaseInstanceID --stack infra)
CURRENT_IP=$(curl -s ifconfig.me | tr -d '\n')

gcloud sql instances patch ${SQL_INSTANCE} --authorized-networks=${CURRENT_IP}

# once you are done with deploying the stacks, you can clear the authorised networks like this
gcloud sql instances patch ${SQL_INSTANCE} --clear-authorized-networks
```

You can now deploy the env-specific infra stack.

```bash
pulumi stack init infra-dev
pulumi stack select infra-dev

API_SERVICE_URL=https://api.example.com pulumi up
```

> **Note:** The `API_SERVICE_URL` is a bit of a chicken or the egg -problem; you need to have the
> env-specific infra stack deployed before deploying the API but the env-specific infra stack needs
> the API endpoint. To get around this, you can first deploy with a random API endpoint, deploy the
> actual API and then re-deploy the env-specific infra stack with the correct endpoint.

> **TODO:** Instructions for deploying the application stacks.

## Stacks

This section describes the general content of each stack. You can find the code from the
`./stack/{stack}.go` files.

### infra

> **TODO:** Describe the `infra` stack's resources.

### infra-env

> **TODO:** Describe the `infra-env` stack's resources.
