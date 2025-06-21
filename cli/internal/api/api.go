package api

import (
	"net/http"
	"net/url"
)

type NitricApiClient struct {
	apiUrl string
}

func NewNitricApiClient(apiUrl string) *NitricApiClient {
	return &NitricApiClient{
		apiUrl: apiUrl,
	}
}

func (c *NitricApiClient) get(path string) (*http.Response, error) {
	apiUrl, err := url.JoinPath(c.apiUrl, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
