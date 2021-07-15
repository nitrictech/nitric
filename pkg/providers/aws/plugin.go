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
	"github.com/nitric-dev/membrane/pkg/plugins/document/dynamodb"
	"github.com/nitric-dev/membrane/pkg/plugins/eventing/sns"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway/lambda"
	"github.com/nitric-dev/membrane/pkg/plugins/queue/sqs"
	"github.com/nitric-dev/membrane/pkg/plugins/storage/s3"
	"github.com/nitric-dev/membrane/pkg/sdk"
)

type AWSServiceFactory struct {
}

func New() sdk.ServiceFactory {
	return &AWSServiceFactory{}
}

// NewDocumentService - Return AWS DynamoDB document plugin
func (p *AWSServiceFactory) NewDocumentService() (sdk.DocumentService, error) {
	return dynamodb_service.New()
}

// NewEventService - Returns AWS SNS based eventing plugin
func (p *AWSServiceFactory) NewEventService() (sdk.EventService, error) {
	return sns_service.New()
}

// NewGatewayService - Returns AWS Lambda Gateway plugin
func (p *AWSServiceFactory) NewGatewayService() (sdk.GatewayService, error) {
	return lambda_service.New()
}

// NewQueueService - Returns AWS SQS based queue plugin
func (p *AWSServiceFactory) NewQueueService() (sdk.QueueService, error) {
	return sqs_service.New()
}

// NewStorageService - Returns AWS S3 based storage plugin
func (p *AWSServiceFactory) NewStorageService() (sdk.StorageService, error) {
	return s3_service.New()
}
