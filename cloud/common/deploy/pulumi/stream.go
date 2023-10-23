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
