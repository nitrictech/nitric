package api

import (
	"net/http"
	"net/url"

	"github.com/nitrictech/nitric/cli/internal/api/transformer"
)

type NitricApiClient struct {
	apiUrl *url.URL
}

func NewNitricApiClient(apiUrl *url.URL) *NitricApiClient {
	return &NitricApiClient{
		apiUrl: apiUrl,
	}
}

func (c *NitricApiClient) get(path string, transformers ...transformer.RequestTransformer) (*http.Response, error) {
	apiUrl, err := url.JoinPath(c.apiUrl.String(), path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	for _, transformer := range transformers {
		err := transformer(req)
		if err != nil {
			return nil, err
		}
	}

	return http.DefaultClient.Do(req)
}
