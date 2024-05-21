package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/queue"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Queue - Deploy a Queue
func (a *NitricAwsTerraformProvider) Queue(stack cdktf.TerraformStack, name string, config *deploymentspb.Queue) error {
	a.Queues[name] = queue.NewQueue(stack, jsii.String(name), &queue.QueueConfig{
		QueueName: jsii.String(name),
		StackId:   a.Stack.StackIdOutput(),
	})

	return nil
}
