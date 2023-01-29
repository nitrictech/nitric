package deploy

import (
	"context"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DownStreamMessageWriter struct {
	stream deploy.DeployService_DownServer
}

func (s *DownStreamMessageWriter) Write(bytes []byte) (int, error) {
	err := s.stream.Send(&deploy.DeployDownEvent{
		Content: &deploy.DeployDownEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: string(bytes),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (d *DeployServer) Down(request *deploy.DeployDownRequest, stream deploy.DeployService_DownServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// TODO: Tear down the requested stack
	dsMessageWriter := &DownStreamMessageWriter{
		stream: stream,
	}

	s, err := auto.UpsertStackInlineSource(context.TODO(), details.Stack, details.Project, nil)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	// destroy the stack
	_, err = s.Destroy(context.TODO(), optdestroy.ProgressStreams(dsMessageWriter))
	if err != nil {
		return err
	}

	return nil
}
