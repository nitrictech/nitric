package main

import (
	cognito_plugin "github.com/nitric-dev/membrane/plugins/auth/cognito"
	sns_plugin "github.com/nitric-dev/membrane/plugins/eventing/sns"
	lambda_plugin "github.com/nitric-dev/membrane/plugins/gateway/lambda"
	dynamodb_plugin "github.com/nitric-dev/membrane/plugins/kv/dynamodb"
	sqs_plugin "github.com/nitric-dev/membrane/plugins/queue/sqs"
	s3_plugin "github.com/nitric-dev/membrane/plugins/storage/s3"
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
