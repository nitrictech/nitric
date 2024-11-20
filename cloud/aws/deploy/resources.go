package deploy

import (
	"encoding/json"

	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAwsPulumiProvider) resourcesStore(ctx *pulumi.Context) error {
	// Build the AWS resource index from the provider information
	// This will be used to store the ARNs of all resources created by the stack
	bucketNameMap := pulumi.StringMap{}
	for name, bucket := range a.Buckets {
		bucketNameMap[name] = bucket.Bucket
	}

	apiArnMap := pulumi.StringMap{}
	for name, api := range a.Apis {
		apiArnMap[name] = api.Arn
	}

	apiEndpointMap := pulumi.StringMap{}
	for name, api := range a.Apis {
		apiEndpointMap[name] = api.ApiEndpoint
	}

	// Build the index from the provider information
	resourceIndexJson := pulumi.All(bucketNameMap, apiArnMap, apiEndpointMap).ApplyT(func(args []interface{}) (string, error) {
		bucketNameMap := args[0].(map[string]string)
		apiArnMap := args[1].(map[string]string)
		apiEndpointMap := args[2].(map[string]string)

		index := common.NewResourceIndex()
		for name, bucket := range bucketNameMap {
			index.Buckets[name] = bucket
		}

		for name, arn := range apiArnMap {
			index.Apis[name] = common.ApiGateway{
				Arn:      arn,
				Endpoint: apiEndpointMap[name],
			}
		}

		indexJson, err := json.Marshal(index)

		if err != nil {
			return "", err
		}

		return string(indexJson), nil
	}).(pulumi.StringOutput)

	_, err := ssm.NewParameter(ctx, "nitric-resource-index", &ssm.ParameterArgs{
		// Create a deterministic name for the resource index
		Name:     pulumi.Sprintf("/nitric/%s/resource-index", a.StackId),
		DataType: pulumi.String("text"),
		Type:     pulumi.String("String"),
		// Store the nitric resource index serialized as a JSON string
		Value: resourceIndexJson,
	})

	// TODO: give services permission to read this parameter store

	return err
}
