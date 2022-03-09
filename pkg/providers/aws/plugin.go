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
	"github.com/nitrictech/nitric/pkg/plugins/document"
	dynamodb_service "github.com/nitrictech/nitric/pkg/plugins/document/dynamodb"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	sns_service "github.com/nitrictech/nitric/pkg/plugins/events/sns"
	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	lambda_service "github.com/nitrictech/nitric/pkg/plugins/gateway/lambda"
	"github.com/nitrictech/nitric/pkg/plugins/queue"
	sqs_service "github.com/nitrictech/nitric/pkg/plugins/queue/sqs"
	"github.com/nitrictech/nitric/pkg/plugins/storage"
	s3_service "github.com/nitrictech/nitric/pkg/plugins/storage/s3"
	"github.com/nitrictech/nitric/pkg/providers"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
)

type AWSServiceFactory struct {
	provider core.AwsProvider
}

func New() (providers.ServiceFactory, error) {
	provider, err := core.New()

	if err != nil {
		return nil, err
	}

	return &AWSServiceFactory{
		provider: provider,
	}, nil
}

// NewDocumentService - Return AWS DynamoDB document plugin
func (p *AWSServiceFactory) NewDocumentService() (document.DocumentService, error) {
	return dynamodb_service.New(p.provider)
}

// NewEventService - Returns AWS SNS based events plugin
func (p *AWSServiceFactory) NewEventService() (events.EventService, error) {
	return sns_service.New(p.provider)
}

// NewGatewayService - Returns AWS Lambda Gateway plugin
func (p *AWSServiceFactory) NewGatewayService() (gateway.GatewayService, error) {
	return lambda_service.New(p.provider)
}

// NewQueueService - Returns AWS SQS based queue plugin
func (p *AWSServiceFactory) NewQueueService() (queue.QueueService, error) {
	return sqs_service.New(p.provider)
}

// NewStorageService - Returns AWS S3 based storage plugin
func (p *AWSServiceFactory) NewStorageService() (storage.StorageService, error) {
	return s3_service.New(p.provider)
}
