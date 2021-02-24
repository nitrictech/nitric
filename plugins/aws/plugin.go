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

type AWSServiceFactory struct {

}

func New() sdk.ServiceFactory {
	return &AWSServiceFactory{}
}

// NewAuthPlugin - Returns AWS Cognito based auth plugin
func (p *AWSServiceFactory) NewAuthService() (sdk.UserService, error) {
	return cognito_plugin.New()
}

// NewDocumentPlugin - Returns AWS DynamoDB based document plugin
func (p *AWSServiceFactory) NewDocumentService() (sdk.DocumentService, error) {
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