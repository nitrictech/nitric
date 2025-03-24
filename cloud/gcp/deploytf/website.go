package deploytf

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/cdn"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/website"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type CdnInput struct {
	ZoneName   *string `json:"zone_name"`
	DomainName *string `json:"domain_name"`
	ClientTtl  *int    `json:"client_ttl"`
	DefaultTtl *int    `json:"default_ttl"`
}

type ApiInput struct {
	Region      *string `json:"region"`
	GatewayId   *string `json:"gateway_id"`
	DefaultHost *string `json:"default_host"`
}

type WebsiteInput struct {
	BasePath       *string `json:"base_path"`
	BucketName     *string `json:"name"`
	IndexDocument  *string `json:"index_document"`
	ErrorDocument  *string `json:"error_document"`
	LocalDirectory *string `json:"local_directory"`
}

func (a *NitricGcpTerraformProvider) deployEntrypoint(stack cdktf.TerraformStack) error {
	if a.GcpConfig.CdnDomain.ZoneName == "" {
		return fmt.Errorf("a valid DNS zone is required to deploy websites to GCP")
	}

	if a.GcpConfig.CdnDomain.DomainName == "" {
		return fmt.Errorf("a valid domain name is required to deploy websites to GCP")
	}

	apis := map[string]ApiInput{}
	websites := map[string]WebsiteInput{}

	cdnInput := &CdnInput{
		ZoneName:   jsii.String(a.GcpConfig.CdnDomain.ZoneName),
		DomainName: jsii.String(a.GcpConfig.CdnDomain.DomainName),
		ClientTtl:  a.GcpConfig.CdnDomain.ClientTtl,
		DefaultTtl: a.GcpConfig.CdnDomain.DefaultTtl,
	}

	for name, api := range a.Apis {
		apis[name] = ApiInput{
			Region:      api.RegionOutput(),
			GatewayId:   api.GatewayIdOutput(),
			DefaultHost: api.DefaultHostOutput(),
		}
	}

	for name, website := range a.Websites {
		websiteName := name
		if *website.BasePath() == "/" {
			websiteName = "default"
		}

		websites[websiteName] = WebsiteInput{
			BasePath:       website.BasePath(),
			BucketName:     website.BucketNameOutput(),
			IndexDocument:  website.IndexDocumentOutput(),
			ErrorDocument:  website.ErrorDocumentOutput(),
			LocalDirectory: website.LocalDirectoryOutput(),
		}
	}

	cdn.NewCdn(stack, jsii.String("cdn"), &cdn.CdnConfig{
		ProjectId:      jsii.String(a.GcpConfig.ProjectId),
		ApiGateways:    apis,
		Region:         jsii.String(a.Region),
		StackId:        a.Stack.StackIdOutput(),
		WebsiteBuckets: websites,
		CdnDomain:      cdnInput,
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
