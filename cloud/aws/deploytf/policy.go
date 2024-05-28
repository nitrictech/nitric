package deploytf

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/policy"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/samber/lo"
)

// func md5Hash(b []byte) string {
// 	hasher := md5.New() //#nosec G401 -- md5 used only to produce a unique ID from non-sensistive information (policy IDs)
// 	hasher.Write(b)

// 	return hex.EncodeToString(hasher.Sum(nil))
// }

var awsActionsMap map[resourcespb.Action][]string = map[resourcespb.Action][]string{
	resourcespb.Action_BucketFileList: {
		"s3:ListBucket",
	},
	resourcespb.Action_BucketFileGet: {
		"s3:GetObject",
	},
	resourcespb.Action_BucketFilePut: {
		"s3:PutObject",
	},
	resourcespb.Action_BucketFileDelete: {
		"s3:DeleteObject",
	},
	resourcespb.Action_TopicPublish: {
		"sns:GetTopicAttributes",
		"sns:Publish",
		"states:StartExecution",
		"states:StateSyncExecution",
	},
	resourcespb.Action_KeyValueStoreRead: {
		"dynamodb:GetItem",
		"dynamodb:BatchGetItem",
		"dynamodb:Scan", // required to scan keys
	},
	resourcespb.Action_KeyValueStoreWrite: {
		"dynamodb:UpdateItem",
		"dynamodb:PutItem",
	},
	resourcespb.Action_KeyValueStoreDelete: {
		"dynamodb:DeleteItem",
	},
	resourcespb.Action_SecretAccess: {
		"secretsmanager:GetSecretValue",
	},
	resourcespb.Action_SecretPut: {
		"secretsmanager:PutSecretValue",
	},
	resourcespb.Action_WebsocketManage: {
		"execute-api:ManageConnections",
	},
	resourcespb.Action_QueueEnqueue: {
		"sqs:SendMessage",
		"sqs:GetQueueAttributes",
		"sqs:GetQueueUrl",
		"sqs:ListQueueTags",
	},
	resourcespb.Action_QueueDequeue: {
		"sqs:ReceiveMessage",
		"sqs:DeleteMessage",
		"sqs:GetQueueAttributes",
		"sqs:GetQueueUrl",
		"sqs:ListQueueTags",
	},
}

func actionsToAwsActions(actions []resourcespb.Action) []string {
	awsActions := make([]string, 0)

	for _, a := range actions {
		awsActions = append(awsActions, awsActionsMap[a]...)
	}

	awsActions = lo.Uniq(awsActions)

	return awsActions
}

// // discover the arn of a deployed resource
func (a *NitricAwsTerraformProvider) arnForResource(resource *deploymentspb.Resource) ([]*string, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Bucket:
		if b, ok := a.Buckets[resource.Id.Name]; ok {
			return []*string{b.BucketArnOutput(), jsii.String(fmt.Sprintf("%s/*", *b.BucketArnOutput()))}, nil
		}
	case resourcespb.ResourceType_Topic:
		if t, ok := a.Topics[resource.Id.Name]; ok {
			return []*string{t.TopicArnOutput(), t.SfnArnOutput()}, nil
		}
	case resourcespb.ResourceType_Queue:
		if q, ok := a.Queues[resource.Id.Name]; ok {
			return []*string{q.QueueArnOutput()}, nil
		}
	case resourcespb.ResourceType_KeyValueStore:
		if c, ok := a.KeyValueStores[resource.Id.Name]; ok {
			return []*string{c.KvArnOutput()}, nil
		}
	case resourcespb.ResourceType_Secret:
		if s, ok := a.Secrets[resource.Id.Name]; ok {
			return []*string{s.SecretArnOutput()}, nil
		}
	case resourcespb.ResourceType_Websocket:
		if w, ok := a.Websockets[resource.Id.Name]; ok {
			return []*string{jsii.String(fmt.Sprintf("%s/*", *w.WebsocketExecArnOutput()))}, nil
		}
	default:
		return nil, fmt.Errorf(
			"invalid resource type: %s. Did you mean to define it as a principal?", resource.Id.Type)
	}

	return nil, fmt.Errorf("unable to find %s named %s in AWS provider resource cache", resource.Id.Type, resource.Id.Name)
}

func (a *NitricAwsTerraformProvider) roleForPrincipal(resource *deploymentspb.Resource) (*string, error) {
	switch resource.Id.Type {
	case resourcespb.ResourceType_Service:
		if f, ok := a.Services[resource.Id.Name]; ok {
			return f.RoleNameOutput(), nil
		}
	default:
		return nil, fmt.Errorf("could not find role for principal: %+v", resource)
	}

	return nil, fmt.Errorf("could not find role for principal: %+v", resource)
}

func (a *NitricAwsTerraformProvider) Policy(stack cdktf.TerraformStack, name string, config *deploymentspb.Policy) error {

	// Get Actions
	actions := actionsToAwsActions(config.Actions)

	// Get Targets
	targetArns := make([]*string, 0, len(config.Resources))

	for _, res := range config.Resources {
		if arn, err := a.arnForResource(res); err == nil {
			targetArns = append(targetArns, arn...)
		} else {
			return fmt.Errorf("failed to create policy, unable to determine resource ARN: %w", err)
		}
	}

	// Get principal roles
	// We're collecting roles here to ensure all defined principals are valid before proceeding
	principalRoles := map[string]*string{}

	for _, princ := range config.Principals {
		if role, err := a.roleForPrincipal(princ); err == nil {
			nameType := fmt.Sprintf("%s:%s", princ.Id.Name, princ.Id.Type)
			if princ.Id.Type != resourcespb.ResourceType_Service {
				return fmt.Errorf("invalid principal type: %s. Only services can be principals", princ.Id.Type)
			}

			principalRoles[nameType] = role
		} else {
			return err
		}
	}

	policy.NewPolicy(stack, jsii.String(name), &policy.PolicyConfig{
		// Name:   jsii.String(name),
		Actions:    jsii.Strings(actions...),
		Resources:  &targetArns,
		Principals: &principalRoles,
	})

	return nil
}
