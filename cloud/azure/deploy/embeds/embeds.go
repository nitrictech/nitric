package embeds

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed api-policy-template.xml
var apiPolicyTemplate string

type ApiPolicyTemplateArgs struct {
	BackendHostName         string
	ExtraPolicies           string
	ManagedIdentityResource string
	ManagedIdentityClientId string
}

func GetApiPolicyTemplate(args ApiPolicyTemplateArgs) (string, error) {
	tmpl, err := template.New("apiPolicyTemplate").Parse(apiPolicyTemplate)
	if err != nil {
		return "", err
	}

	var output strings.Builder

	err = tmpl.Execute(&output, args)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

type JwtTemplateArgs struct {
	OidcUri       string
	RequiredClaim string
}

//go:embed api-jwt-template.xml
var jwtPolicyTemplate string

func GetApiJwtTemplate(args JwtTemplateArgs) (string, error) {
	tmpl, err := template.New("jwtPolicyTemplate").Parse(jwtPolicyTemplate)
	if err != nil {
		return "", err
	}

	var output strings.Builder

	err = tmpl.Execute(&output, args)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
