package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/config"
	detail "github.com/nitrictech/nitric/cli/internal/details"
	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/samber/do/v2"
)

type AuthDetails struct {
	WorkOS detail.WorkOSDetails `json:"workos"`
}

type Service struct {
	nitricBackendUrl *url.URL
}

var _ detail.AuthDetailsService = &Service{}

func NewService(inj do.Injector) (*Service, error) {
	conf := do.MustInvoke[*config.Config](inj)

	return &Service{nitricBackendUrl: conf.GetNitricServerUrl()}, nil
}

func (s *Service) GetWorkOSDetails() (*detail.WorkOSDetails, error) {
	apiUrl, err := url.JoinPath(s.nitricBackendUrl.String(), "/auth/details")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset by peer") {
			return nil, fmt.Errorf("failed to connect to the %s API. Please check your connection and try again. If the problem persists, please contact support.", version.ProductName)
		}
		return nil, fmt.Errorf("failed to connect to %s auth details endpoint: %v", version.ProductName, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s auth details endpoint: %v", version.ProductName, err)
	}

	var authDetails AuthDetails
	err = json.Unmarshal(body, &authDetails)
	if err != nil {
		return nil, fmt.Errorf("unexpected response from %s auth details endpoint: %v", version.ProductName, err)
	}

	return &authDetails.WorkOS, nil
}
