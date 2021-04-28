// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
