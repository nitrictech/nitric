// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsbatch "github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/google/uuid"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type AwsBatchService struct {
	stackId     string
	client      *awsbatch.Client
	jobQueueArn string
	batchpb.UnimplementedBatchServer
}

func (a *AwsBatchService) SubmitJob(ctx context.Context, request *batchpb.JobSubmitRequest) (*batchpb.JobSubmitResponse, error) {
	// find and submit a new job
	jobDefinitionName, err := common.GetJobDefinitionName(a.stackId, request.GetJobName())
	if err != nil {
		return nil, err
	}

	jobName := uuid.New()

	fmt.Printf("Submitting job to AWS Batch for JD: %s onto: %s \n", jobDefinitionName, a.jobQueueArn)

	jsonData, err := protojson.Marshal(request.GetData())
	if err != nil {
		return nil, err
	}

	_, err = a.client.SubmitJob(ctx, &awsbatch.SubmitJobInput{
		JobDefinition: aws.String(jobDefinitionName),
		// Generate a unique name for the job
		JobName:  aws.String(fmt.Sprintf("%s-%s", jobName, request.GetJobName())),
		JobQueue: aws.String(a.jobQueueArn),
		ContainerOverrides: &types.ContainerOverrides{
			Environment: []types.KeyValuePair{
				{
					Name:  aws.String("NITRIC_JOB_DATA"),
					Value: aws.String(string(jsonData)),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("Job submitted to AWS Batch")

	return &batchpb.JobSubmitResponse{}, nil
}

func New() (*AwsBatchService, error) {
	jobQueueArn := env.JOB_QUEUE_ARN.String()
	stackId := commonenv.NITRIC_STACK_ID.String()

	if jobQueueArn == "" {
		return nil, status.Error(codes.InvalidArgument, "JOB_QUEUE_ARN not set")
	}

	awsRegion := env.AWS_REGION.String()

	// Create a new AWS session
	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion), config.WithClientLogMode(aws.LogRetries|aws.LogResponse|aws.LogRequest))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	return &AwsBatchService{
		// TODO: Configure client
		stackId:     stackId,
		client:      awsbatch.NewFromConfig(cfg),
		jobQueueArn: jobQueueArn,
	}, nil
}
