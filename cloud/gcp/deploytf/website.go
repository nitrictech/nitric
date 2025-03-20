package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/website"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricGcpTerraformProvider) Website(stack cdktf.TerraformStack, name string, config *deploymentspb.Website) error {
	// Deploy a website
	a.Websites[name] = website.NewWebsite(stack, jsii.Sprintf("website_%s", name), &website.WebsiteConfig{
		WebsiteName:    jsii.String(name),
		StackId:        a.Stack.StackIdOutput(),
		BasePath:       jsii.String(config.BasePath),
		LocalDirectory: jsii.String(config.GetLocalDirectory()),
		Region:         jsii.String(a.Region),
		ErrorDocument:  jsii.String(config.GetErrorDocument()),
		IndexDocument:  jsii.String(config.GetIndexDocument()),
	})

	return nil
}
