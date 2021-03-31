package main

import (
	cognito_plugin "github.com/nitric-dev/membrane/plugins/aws/auth/cognito"
	dynamodb_plugin "github.com/nitric-dev/membrane/plugins/aws/kv/dynamodb"
	sns_plugin "github.com/nitric-dev/membrane/plugins/aws/eventing/sns"
	lambda_plugin "github.com/nitric-dev/membrane/plugins/aws/gateway/lambda"
	sqs_plugin "github.com/nitric-dev/membrane/plugins/aws/queue/sqs"
	s3_plugin "github.com/nitric-dev/membrane/plugins/aws/storage/s3"
	"github.com/nitric-dev/membrane/sdk"
)

type AWSServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &AWSServiceFactory{}
}

// NewAuthPlugin - Returns AWS Cognito based auth plugin
func (p *AWSServiceFactory) NewAuthService() (sdk.UserService, error) {
	return cognito_plugin.New()
}

// NewDocumentPlugin - Returns AWS DynamoDB based key value plugin
func (p *AWSServiceFactory) NewKeyValueService() (sdk.KeyValueService, error) {
	return dynamodb_plugin.New()
}

// NewEventingPlugin - Returns AWS SNS based eventing plugin
func (p *AWSServiceFactory) NewEventService() (sdk.EventService, error) {
	return sns_plugin.New()
}

// NewGatewayPlugin - Returns AWS Lambda Gateway plugin
func (p *AWSServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return lambda_plugin.New()
}

// NewQueuePlugin - Returns AWS SQS based queue plugin
func (p *AWSServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return sqs_plugin.New()
}

// NewStoragePlugin - Returns AWS S3 based storage plugin
func (p *AWSServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return s3_plugin.New()
}
