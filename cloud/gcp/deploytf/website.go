package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/cdn"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/website"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type ApiInput struct {
	Region      *string `json:"region"`
	GatewayId   *string `json:"gateway_id"`
	DefaultHost *string `json:"default_host"`
}

type WebsiteInput struct {
	BucketName     *string `json:"bucket_name"`
	IndexDocument  *string `json:"base_path"`
	ErrorDocument  *string `json:"error_document"`
	LocalDirectory *string `json:"local_directory"`
}

func (a *NitricGcpTerraformProvider) deployEntrypoint(stack cdktf.TerraformStack) error {
	apis := map[string]ApiInput{}
	websites := map[string]WebsiteInput{}

	for name, api := range a.Apis {
		apis[name] = ApiInput{
			Region:      api.RegionOutput(),
			GatewayId:   api.GatewayIdOutput(),
			DefaultHost: api.DefaultHostOutput(),
		}
	}

	for name, website := range a.Websites {
		websites[name] = WebsiteInput{
			BucketName:     website.BucketNameOutput(),
			IndexDocument:  website.IndexDocumentOutput(),
			ErrorDocument:  website.ErrorDocumentOutput(),
			LocalDirectory: website.LocalDirectoryOutput(),
		}
	}

	cdn.NewCdn(stack, jsii.String("cdn"), &cdn.CdnConfig{
		ApiGateways:    apis,
		Region:         jsii.String(a.Region),
		StackId:        a.Stack.StackIdOutput(),
		WebsiteBuckets: websites,
		CdnDomain:      &a.GcpConfig.CdnDomain,
	})

	return nil
}

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
