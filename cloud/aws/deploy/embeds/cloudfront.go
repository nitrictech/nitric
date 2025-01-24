package embeds

import (
	_ "embed"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed api-url-rewrite.js
var cloudfront_ApiUrlRewriteFunction string

//go:embed url-rewrite.js
var cloudfront_UrlRewriteFunctionName string

func GetApiUrlRewriteFunction() pulumi.StringInput {
	return pulumi.String(cloudfront_ApiUrlRewriteFunction)
}

func GetUrlRewriteFunction() pulumi.StringInput {
	return pulumi.String(cloudfront_UrlRewriteFunctionName)
}
