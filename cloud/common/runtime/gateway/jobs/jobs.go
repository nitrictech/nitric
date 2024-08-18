package jobs

import (
	"fmt"
	"log"

	"github.com/nitrictech/nitric/cloud/common/runtime/env"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type DefaultBatchGateway struct {
	gateway.UnimplementedGatewayPlugin
}

func (s *DefaultBatchGateway) Start(opts *gateway.GatewayStartOpts) error {
	// all of our workers should be available now to process jobs

	jobName := env.NITRIC_JOB_NAME.String()
	jobData := env.NITRIC_JOB_DATA.String()

	jobDataProto := &batchpb.JobData{}

	err := protojson.Unmarshal([]byte(jobData), jobDataProto)
	if err != nil {
		return fmt.Errorf("unable to unmarshal job data: %v", err)
	}

	// construct the job event
	response, err := opts.JobHandlerPlugin.HandleJobRequest(&batchpb.ServerMessage{
		Content: &batchpb.ServerMessage_JobRequest{
			JobRequest: &batchpb.JobRequest{
				JobName: jobName,
				Data:    jobDataProto,
			},
		},
	})

	if err != nil || !response.GetJobResponse().Success {
		log.Fatalf("Job failed to successfully execute: %v", err)
	}

	return nil
}

func (s *DefaultBatchGateway) Stop() error {
	// No-op, all work is completed as part of the gateway start
	// the gateway simply blocks until the job has been processed
	return nil
}

func NewDefaultBatchGateway() *DefaultBatchGateway {
	return &DefaultBatchGateway{}
}
