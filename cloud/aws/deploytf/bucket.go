package deploytf

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/bucket"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// Bucket - Deploy a Storage Bucket
func (n *NitricAwsTerraformProvider) Bucket(stack cdktf.TerraformStack, name string, config *deploymentspb.Bucket) error {
	n.Buckets[name] = bucket.NewBucket(stack, &name, &bucket.BucketConfig{
		BucketName: &name,
		StackId:    n.Stack.StackIdOutput(),
	})

	// TODO: Deploy bucket subscriptions

	return nil
}
