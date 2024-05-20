package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/topic"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAwsTerraformProvider) Topic(stack cdktf.TerraformStack, name string, config *deploymentspb.Topic) error {
	lambdaSubscriberArns := []*string{}

	for _, subscriber := range config.Subscriptions {
		// subscriber.GetService()
		lambdaService := a.Services[subscriber.GetService()]
		lambdaSubscriberArns = append(lambdaSubscriberArns, lambdaService.LambdaArnOutput())
	}

	a.Topics[name] = topic.NewTopic(stack, &name, &topic.TopicConfig{
		StackId:           a.Stack.StackIdOutput(),
		TopicName:         jsii.String(name),
		LambdaSubscribers: &lambdaSubscriberArns,
	})

	return nil
}
