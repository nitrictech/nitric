package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nitrictech/nitric/cli/internal/config"
	"github.com/pkg/errors"
	"github.com/samber/do/v2"
)

type TokenProvider interface {
	// GetAccessToken returns the access token for the user
	GetAccessToken() (string, error)
}

type NitricApiClient struct {
	tokenProvider TokenProvider
	apiUrl        *url.URL
}

func NewNitricApiClient(injector do.Injector) (*NitricApiClient, error) {
	config := do.MustInvoke[*config.Config](injector)
	apiUrl := config.GetNitricServerUrl()

	tokenProvider := do.MustInvokeAs[TokenProvider](injector)

	return &NitricApiClient{
		apiUrl:        apiUrl,
		tokenProvider: tokenProvider,
	}, nil
}

func (c *NitricApiClient) get(path string, requiresAuth bool) (*http.Response, error) {
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

		token, err := c.tokenProvider.GetAccessToken()
		if err != nil {
			return nil, errors.Wrap(ErrUnauthenticated, err.Error())
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return http.DefaultClient.Do(req)
}
