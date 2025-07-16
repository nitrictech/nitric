package devserver

import (
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/build"
)

type NitricProjectBuild struct {
	apiClient *api.NitricApiClient
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

func (n *NitricProjectBuild) OnConnect(send SendFunc) {
	// No-op
}

func (n *NitricProjectBuild) OnMessage(message json.RawMessage) {
	var buildMessage Message[ProjectBuild]
	err := json.Unmarshal(message, &buildMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message", err)
		return
	}

	// Not the right message type continue
	if buildMessage.Type != "nitricBuild" {
		return
	}

	stackPath, err := n.builder.BuildProjectFromFileForTarget("nitric.yaml", buildMessage.Payload.Target)
	if err != nil {
		n.broadcast(Message[any]{
			Type: "nitricBuildError",
			Payload: ProjectBuildError{
				Message: err.Error(),
			},
		})
		return
	}

	n.broadcast(Message[any]{
		Type: "nitricBuildSuccess",
		Payload: ProjectBuildSuccess{
			StackPath: stackPath,
		},
	})
}

func NewProjectBuild(apiClient *api.NitricApiClient, builder *build.BuilderService, broadcast BroadcastFunc) (*NitricProjectBuild, error) {
	buildServer := &NitricProjectBuild{
		apiClient: apiClient,
		broadcast: broadcast,
		builder:   builder,
	}

	return buildServer, nil
}
