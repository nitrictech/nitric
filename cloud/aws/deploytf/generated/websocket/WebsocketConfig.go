package websocket

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type WebsocketConfig struct {
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Providers *[]interface{} `field:"optional" json:"providers" yaml:"providers"`
	// Experimental.
	SkipAssetCreationFromLocalModules *bool `field:"optional" json:"skipAssetCreationFromLocalModules" yaml:"skipAssetCreationFromLocalModules"`
	// The ARN of the lambda to send websocket connection events to.
	LambdaConnectTarget *string `field:"required" json:"lambdaConnectTarget" yaml:"lambdaConnectTarget"`
	// The ARN of the lambda to send websocket disconnection events to.
	LambdaDisconnectTarget *string `field:"required" json:"lambdaDisconnectTarget" yaml:"lambdaDisconnectTarget"`
	// The ARN of the lambda to send websocket disconnection events to.
	LambdaMessageTarget *string `field:"required" json:"lambdaMessageTarget" yaml:"lambdaMessageTarget"`
	// The ID of the Nitric stack.
	StackId *string `field:"required" json:"stackId" yaml:"stackId"`
	// The name of the websocket.
	WebsocketName *string `field:"required" json:"websocketName" yaml:"websocketName"`
}
