package devserver

import (
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/build"
	"github.com/nitrictech/nitric/cli/internal/version"
)

type SugaProjectBuild struct {
	apiClient *api.SugaApiClient
	broadcast BroadcastFunc
	builder   *build.BuilderService
}

type ProjectBuild struct {
	Target string `json:"target"`
}

type ProjectBuildSuccess struct {
	StackPath string `json:"stackPath"`
}

type ProjectBuildError struct {
	Message string `json:"message"`
}

func (n *SugaProjectBuild) OnConnect(send SendFunc) {
	// No-op
}

func (n *SugaProjectBuild) OnMessage(message json.RawMessage) {
	var buildMessage Message[ProjectBuild]
	err := json.Unmarshal(message, &buildMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message", err)
		return
	}

	// Not the right message type continue
	if buildMessage.Type != "buildMessage" {
		return
	}

	stackPath, err := n.builder.BuildProjectFromFileForTarget(version.ConfigFileName, buildMessage.Payload.Target)
	if err != nil {
		n.broadcast(Message[any]{
			Type: "buildError",
			Payload: ProjectBuildError{
				Message: err.Error(),
			},
		})
		return
	}

	n.broadcast(Message[any]{
		Type: "buildSuccess",
		Payload: ProjectBuildSuccess{
			StackPath: stackPath,
		},
	})
}

func NewProjectBuild(apiClient *api.SugaApiClient, builder *build.BuilderService, broadcast BroadcastFunc) (*SugaProjectBuild, error) {
	buildServer := &SugaProjectBuild{
		apiClient: apiClient,
		broadcast: broadcast,
		builder:   builder,
	}

	return buildServer, nil
}
