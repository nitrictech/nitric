package main

import (
	cognito_plugin "github.com/nitric-dev/membrane/plugins/aws/auth/cognito"
	dynamodb_plugin "github.com/nitric-dev/membrane/plugins/aws/documents/dynamodb"
	sns_plugin "github.com/nitric-dev/membrane/plugins/aws/eventing/sns"
	lambda_plugin "github.com/nitric-dev/membrane/plugins/aws/gateway/lambda"
	sqs_plugin "github.com/nitric-dev/membrane/plugins/aws/queue/sqs"
	s3_plugin "github.com/nitric-dev/membrane/plugins/aws/storage/s3"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

type AWSPluginFactory struct {

}

func New() sdk.ServiceFactory {
	return &AWSPluginFactory{}
}

// NewAuthPlugin - Returns AWS Cognito based auth plugin
func (p *AWSPluginFactory) NewAuthService() (sdk.AuthService, error) {
	return cognito_plugin.New()
}

// NewDocumentPlugin - Returns AWS DynamoDB based document plugin
func (p *AWSPluginFactory) NewDocumentService() (sdk.DocumentService, error) {
	return dynamodb_plugin.New()
}

// NewEventingPlugin - Returns AWS SNS based eventing plugin
func (p *AWSPluginFactory) NewEventService() (sdk.EventService, error) {
	return sns_plugin.New()
}

// NewGatewayPlugin - Returns AWS Lambda Gateway plugin
func (p *AWSPluginFactory) NewGatewayService() (sdk.GatewayService, error) {
	return lambda_plugin.New()
}

// NewQueuePlugin - Returns AWS SQS based queue plugin
func (p *AWSPluginFactory) NewQueueService() (sdk.QueueService, error) {
	return sqs_plugin.New()
}

// NewStoragePlugin - Returns AWS S3 based storage plugin
func (p *AWSPluginFactory) NewStorageService() (sdk.StorageService, error) {
	return s3_plugin.New()
}