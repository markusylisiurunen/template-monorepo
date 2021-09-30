package stack

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Stack interface {
	Run(ctx *pulumi.Context) error
}
