package pulumi

import (
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type UpStreamMessageWriter struct {
	Stream deploy.DeployService_UpServer
}

func (s *UpStreamMessageWriter) Write(bytes []byte) (int, error) {
	str := string(bytes)

	if str == "." {
		// skip progress dots
		return len(bytes), nil
	}

	err := s.Stream.Send(&deploy.DeployUpEvent{
		Content: &deploy.DeployUpEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: str,
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

type DownStreamMessageWriter struct {
	Stream deploy.DeployService_DownServer
}

func (s *DownStreamMessageWriter) Write(bytes []byte) (int, error) {
	str := string(bytes)

	if str == "." {
		// skip progress dots
		return len(bytes), nil
	}

	err := s.Stream.Send(&deploy.DeployDownEvent{
		Content: &deploy.DeployDownEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: str,
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}
