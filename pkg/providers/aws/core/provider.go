package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"

	"github.com/nitrictech/nitric/pkg/utils"
)

type AwsResource = string

const (
	AwsResource_Topic      AwsResource = "sns:topic"
	AwsResource_Collection AwsResource = "dynamodb:table"
	AwsResource_Queue      AwsResource = "sqs:queue"
	AwsResource_Bucket     AwsResource = "s3:bucket"
	AwsResource_Secret     AwsResource = "secretsmanager:secret"
)

type AwsProvider interface {
	// GetResources API operation for AWS Provider.
	// Returns requested aws resources for the given resource type
	GetResources(AwsResource) (map[string]string, error)
}

// Aws core utility provider
type awsProviderImpl struct {
	stack  string
	client *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	cache  map[AwsResource]map[string]string
}

var _ AwsProvider = &awsProviderImpl{}

func (a *awsProviderImpl) GetResources(typ AwsResource) (map[string]string, error) {
	if a.cache[typ] == nil {
		resources := make(map[string]string)

		out, err := a.client.GetResources(&resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []*string{aws.String(typ)},
			TagFilters: []*resourcegroupstaggingapi.TagFilter{{
				Key:    aws.String("x-nitric-stack"),
				Values: []*string{aws.String(a.stack)},
			}, {
				Key: aws.String("x-nitric-name"),
			}},
		})

		if err != nil {
			return nil, err
		}

		for _, tm := range out.ResourceTagMappingList {
			for _, t := range tm.Tags {
				if *t.Key == "x-nitric-name" {
					resources[*t.Value] = *tm.ResourceARN
					break
				}
			}
		}

		a.cache[typ] = resources
	}

	return a.cache[typ], nil
}

func New() (AwsProvider, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, err
	}

	client := resourcegroupstaggingapi.New(sess)

	return &awsProviderImpl{
		client: client,
		cache:  make(map[AwsResource]map[string]string),
	}, nil
}
