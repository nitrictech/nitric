package api

import (
	"net/http"
	"net/url"

	"github.com/nitrictech/nitric/cli/internal/api/transformer"
)

type NitricApiClient struct {
	apiUrl       *url.URL
	transformers []transformer.RequestTransformer
}

func withAcceptHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
}

func NewNitricApiClient(apiUrl *url.URL, transformers ...transformer.RequestTransformer) *NitricApiClient {
	defaultTransformers := []transformer.RequestTransformer{
		withAcceptHeader,
	}

	return &NitricApiClient{
		apiUrl:       apiUrl,
		transformers: append(defaultTransformers, transformers...),
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

	for _, transformer := range c.transformers {
		transformer(req)
	}

	return http.DefaultClient.Do(req)
}
