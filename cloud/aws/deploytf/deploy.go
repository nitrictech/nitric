package deploytf

import (
	"embed"
	"io/fs"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/schedule"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/secret"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	tfstack "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/topic"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsTerraformProvider struct {
	*deploy.CommonStackDetails
	Stack tfstack.Stack

	AwsConfig *common.AwsConfig
	Apis      map[string]api.Api
	Buckets   map[string]bucket.Bucket
	Topics    map[string]topic.Topic
	Schedules map[string]schedule.Schedule
	Services  map[string]service.Service
	Secrets   map[string]secret.Secret
	Queues    map[string]queue.Queue
	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricAwsTerraformProvider)(nil)

func (a *NitricAwsTerraformProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	a.AwsConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

// embed the modules directory here
//
//go:embed modules/**/*
var modules embed.FS

func (a *NitricAwsTerraformProvider) CdkTfModules() (fs.FS, error) {
	return modules, nil
}

func (a *NitricAwsTerraformProvider) Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error {
	a.Stack = tfstack.NewStack(stack, jsii.String("stack"), &tfstack.StackConfig{})

	return nil
}

func (a *NitricAwsTerraformProvider) Post(stack cdktf.TerraformStack) error {
	return nil
}

// // Post - Called after all resources have been created, before the Pulumi Context is concluded
// Post(stack cdktf.TerraformStack) error

func NewNitricAwsProvider() *NitricAwsTerraformProvider {
	return &NitricAwsTerraformProvider{
		Apis:      make(map[string]api.Api),
		Buckets:   make(map[string]bucket.Bucket),
		Services:  make(map[string]service.Service),
		Topics:    make(map[string]topic.Topic),
		Schedules: make(map[string]schedule.Schedule),
		Secrets:   make(map[string]secret.Secret),
		Queues:    make(map[string]queue.Queue),
	}
}
