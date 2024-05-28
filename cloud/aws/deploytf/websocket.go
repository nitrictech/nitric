package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/websocket"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAwsTerraformProvider) Websocket(stack cdktf.TerraformStack, name string, config *deploymentspb.Websocket) error {
	connectTarget := a.Services[config.ConnectTarget.GetService()]
	messageTarget := a.Services[config.MessageTarget.GetService()]
	disconnectTarget := a.Services[config.DisconnectTarget.GetService()]

	a.Websockets[name] = websocket.NewWebsocket(stack, jsii.String(name), &websocket.WebsocketConfig{
		StackId:                a.Stack.StackIdOutput(),
		WebsocketName:          jsii.String(name),
		LambdaConnectTarget:    connectTarget.LambdaArnOutput(),
		LambdaMessageTarget:    messageTarget.LambdaArnOutput(),
		LambdaDisconnectTarget: disconnectTarget.LambdaArnOutput(),
	})

	return nil
}
