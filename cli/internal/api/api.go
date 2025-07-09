package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type TokenProvider interface {
	// GetAccessToken returns the access token for the user
	GetAccessToken() (string, error)
}

type NitricApiClient struct {
	tokenProvider TokenProvider
	apiUrl        *url.URL
}

func NewNitricApiClient(apiUrl *url.URL, tokenProvider TokenProvider) *NitricApiClient {
	return &NitricApiClient{
		apiUrl:        apiUrl,
		tokenProvider: tokenProvider,
	}
}

func (c *NitricApiClient) SetTokenProvider(tokenProvider TokenProvider) {
	c.tokenProvider = tokenProvider
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
