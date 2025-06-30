package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type NitricApiClient struct {
	apiUrl      *url.URL
	client      *http.Client
	accessToken *string
}

func NewNitricApiClient(apiUrl *url.URL, accessToken *string) *NitricApiClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &NitricApiClient{
		apiUrl:      apiUrl,
		client:      client,
		accessToken: accessToken,
	}
}

func (c *NitricApiClient) get(path string) (*http.Response, error) {
	apiUrl, err := url.JoinPath(c.apiUrl.String(), path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	if c.accessToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *c.accessToken))
	}

	req.Header.Set("User-Agent", "nitric-cli")
	req.Header.Set("Accept", "application/json")

	return c.client.Do(req)
}
