package devserver

import (
	"encoding/json"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/cli/internal/platforms"
	"github.com/nitrictech/nitric/cli/internal/plugins"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines/terraform"
	"github.com/spf13/afero"
)

type NitricProjectBuild struct {
	fs        afero.Fs
	apiClient *api.NitricApiClient
	broadcast BroadcastFunc
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

	// If we're sent a build command then start building the project
	appSpec, err := schema.LoadFromFile(n.fs, "nitric.yaml", true)
	if err != nil {
		// TODO: Log/Broadcast the error
		n.broadcast(Message[any]{
			Type: "nitricBuildError",
			Payload: ProjectBuildError{
				Message: err.Error(),
			},
		})
	}

	platformRepository := platforms.NewPlatformRepository(n.apiClient)

	// TODO: Do we care if the file contains no targets if the target we were given was in the message
	// Probably but can add this error case back in if we want to
	// if len(appSpec.Targets) == 0 {
	// 	n.broadcast(Message[any]{
	// 		Type: "nitricBuildError",
	// 		Payload: ProjectBuildError{
	// 			Message: "No targets specified in nitric.yaml",
	// 		},
	// 	})
	// 	return
	// }

	platform, err := terraform.PlatformFromId(n.fs, buildMessage.Payload.Target, platformRepository)
	if err != nil {
		n.broadcast(Message[any]{
			Type: "nitricBuildError",
			Payload: ProjectBuildError{
				Message: err.Error(),
			},
		})
		return
	}

	repo := plugins.NewPluginRepository(n.apiClient)
	engine := terraform.New(platform, terraform.WithRepository(repo))

	stackPath, err := engine.Apply(appSpec)
	if err != nil {
		fmt.Print("Error applying platform: ", err)
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

func NewProjectBuild(fs afero.Fs, apiClient *api.NitricApiClient, broadcast BroadcastFunc) (*NitricProjectBuild, error) {
	buildServer := &NitricProjectBuild{
		fs:        fs,
		apiClient: apiClient,
		broadcast: broadcast,
	}

	return buildServer, nil
}
