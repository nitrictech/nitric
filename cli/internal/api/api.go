package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/pkg/errors"
	"github.com/samber/do/v2"
)

type TokenProvider interface {
	// GetAccessToken returns the access token for the user
	GetAccessToken(forceRefresh bool) (string, error)
}

type SugaApiClient struct {
	tokenProvider TokenProvider
	apiUrl        *url.URL
}

func NewSugaApiClient(injector do.Injector) (*SugaApiClient, error) {
	config, err := do.Invoke[*config.Config](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	apiUrl := config.GetNitricServerUrl()

	tokenProvider, err := do.InvokeAs[TokenProvider](injector)
	if err != nil {
		return nil, fmt.Errorf("failed to get token provider: %w", err)
	}

	return &SugaApiClient{
		apiUrl:        apiUrl,
		tokenProvider: tokenProvider,
	}, nil
}

func (c *SugaApiClient) get(path string, requiresAuth bool) (*http.Response, error) {
	apiUrl, err := url.JoinPath(c.apiUrl.String(), path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	if requiresAuth {
		if c.tokenProvider == nil {
			return nil, errors.Wrap(ErrPreconditionFailed, "no token provider provided")
		}

		token, err := c.tokenProvider.GetAccessToken(false)
		if err != nil {
			return nil, errors.Wrap(ErrUnauthenticated, err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return http.DefaultClient.Do(req)
}

func (c *SugaApiClient) post(path string, requiresAuth bool, body []byte) (*http.Response, error) {
	apiUrl, err := url.JoinPath(c.apiUrl.String(), path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if requiresAuth {
		if c.tokenProvider == nil {
			return nil, errors.Wrap(ErrPreconditionFailed, "no token provider provided")
		}

		token, err := c.tokenProvider.GetAccessToken(false)
		if err != nil {
			return nil, errors.Wrap(ErrUnauthenticated, err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return http.DefaultClient.Do(req)
}
