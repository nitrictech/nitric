package bucket

import (
	"cloud.google.com/go/storage"
	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)


func EventTypeToStorageEventType(eventType *v1.EventType) []string {
	switch *eventType {
	case v1.EventType_All:
		return []string{"OBJECT_FINALIZE", "OBJECT_DELETE"}
	case v1.EventType_Created:
		return []string{"OBJECT_FINALIZE"}
	case v1.EventType_Deleted:
		return []string{"OBJECT_DELETE"}
	default:
		return []string{}
	}
}

type S3Notification struct {
	pulumi.ResourceState

	Name         string
	Notification *storage.Notification
}

type S3NotificationArgs struct {
	Location  string
	StackID   pulumi.StringInput

	Bucket *S3Bucket
	Config *v1.BucketNotificationConfig
	Function *exec.LambdaExecUnit
}

func NewS3Notification(ctx *pulumi.Context, name string, args *S3NotificationArgs, opts ...pulumi.ResourceOption) (*S3Notification, error) {
	res := &S3Notification{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:bucket:GCPCloudStorageNotification", name, res, opts...)
	if err != nil {
		return nil, err
	}

	topicPolicyDocument := iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
		Statements: iam.GetPolicyDocumentStatementArray{
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Principals: iam.GetPolicyDocumentStatementPrincipalArray{
					&iam.GetPolicyDocumentStatementPrincipalArgs{
						Type: pulumi.String("Service"),
						Identifiers: pulumi.StringArray{
							pulumi.String("s3.amazonaws.com"),
						},
					},
				},
				Actions: pulumi.StringArray{
					pulumi.String("SNS:Publish"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("arn:aws:sns:*:*:s3-event-notification-topic"),
				},
				Conditions: iam.GetPolicyDocumentStatementConditionArray{
					&iam.GetPolicyDocumentStatementConditionArgs{
						Test:     pulumi.String("ArnLike"),
						Variable: pulumi.String("aws:SourceArn"),
						Values: pulumi.StringArray{
							args.Bucket.S3.Arn,
						},
					},
				},
			},
		},
	}, nil)

	topic, err := sns.NewTopic(ctx, name+"-topic", &sns.TopicArgs{
		Policy: topicPolicyDocument.ApplyT(func(topicPolicyDocument iam.GetPolicyDocumentResult) (*string, error) {
			return &topicPolicyDocument.Json, nil
		}).(pulumi.StringPtrOutput),
	})
	if err != nil {
		return nil, err
	}

	_, err = awslambda.NewPermission(ctx, name+"-permission", &awslambda.PermissionArgs{
		SourceArn: topic.Arn,
		Function:  args.Function.Function.Name,
		Principal: pulumi.String("sns.amazonaws.com"),
		Action:    pulumi.String("lambda:InvokeFunction"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	snsSubscription, err = sns.NewTopicSubscription(ctx, name+"-sub", &sns.TopicSubscriptionArgs{
		Endpoint: "",
		Protocol: pulumi.String("http"),
		Topic:    topic.ID(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	if args.Config.EventFilter == "*" {
		args.Config.EventFilter = ""
	}

	_, err = s3.NewBucketNotification(ctx, name, &s3.BucketNotificationArgs{
		Bucket: args.Bucket.S3.ID(),
		Topics: s3.BucketNotificationTopicArray{
			&s3.BucketNotificationTopicArgs{
				TopicArn: topic.Arn,
				Events: pulumi.StringArray{
					pulumi.String(""),
				},
				FilterPrefix: pulumi.String(args.Config.EventFilter),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}