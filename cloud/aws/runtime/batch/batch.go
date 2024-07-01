package batch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsbatch "github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/google/uuid"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AwsBatchService struct {
	stackId     string
	client      *awsbatch.Client
	jobQueueArn string
	batchpb.UnimplementedBatchServer
}

func (a *AwsBatchService) CreateJob(ctx context.Context, request *batchpb.CreateJobRequest) (*batchpb.CreateJobResponse, error) {
	// find and submit a new job
	jobDefinitionName, err := common.GetJobDefinitionName(a.stackId, request.GetName())
	if err != nil {
		return nil, err
	}

	jobName := uuid.New()

	_, err = a.client.SubmitJob(ctx, &awsbatch.SubmitJobInput{
		JobDefinition: aws.String(jobDefinitionName),
		// Generate a unique name for the job
		JobName:  aws.String(fmt.Sprintf("%s-%s", jobName, request.GetName())),
		JobQueue: aws.String(a.jobQueueArn),
	})

	if err != nil {
		return nil, err
	}

	return &batchpb.CreateJobResponse{}, nil
}

func New() (*AwsBatchService, error) {
	jobQueueArn := env.JOB_QUEUE_ARN.String()
	stackId := commonenv.NITRIC_STACK_ID.String()

	if jobQueueArn == "" {
		return nil, status.Error(codes.InvalidArgument, "JOB_QUEUE_ARN not set")
	}

	return &AwsBatchService{
		// TODO: Configure client
		stackId:     stackId,
		client:      awsbatch.New(awsbatch.Options{}),
		jobQueueArn: jobQueueArn,
	}, nil
}
