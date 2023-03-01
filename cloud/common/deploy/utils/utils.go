package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func IntValueOrDefault(v, def int) int {
	if v != 0 {
		return v
	}

	return def
}

func StringTrunc(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max]
}

type NamedResource struct {
	Name string
}

func MapNamedResource(resource []*NamedResource) map[string]*NamedResource {
	var resources map[string]*NamedResource
	for _, r := range resource {
		resources[r.Name] = r
	}
	return resources
}

type OpenIdConfig struct {
	Issuer        string `json:"issuer"`
	JwksUri       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
	AuthEndpoint  string `json:"authorization_endpoint"`
}

func GetOpenIdConnectConfig(openIdConnectUrl string) (*OpenIdConfig, error) {
	// append well-known configuration to issuer
	url, err := url.Parse(openIdConnectUrl)
	if err != nil {
		return nil, err
	}

	// get the configuration document
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 status retrieving openid-configuration: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	oidConf := &OpenIdConfig{}

	if err := json.Unmarshal(body, oidConf); err != nil {
		return nil, errors.WithMessage(err, "error unmarshalling open id config")
	}

	return oidConf, nil
}

func GetAudiencesFromExtension(extensions map[string]interface{}) ([]string, error){
	audExt, ok := extensions["x-nitric-audiences"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to get audiences from api spec")
	}

	audiences := make([]string, len(audExt))
	for i, v := range audExt {
		audiences[i] = fmt.Sprint(v)
	}

	return audiences, nil
}