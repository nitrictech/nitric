// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func GetAudiencesFromExtension(extensions map[string]interface{}) ([]string, error) {
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
