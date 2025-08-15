package service

import (
	"context"
	"net/http"
)

type Proxy interface {
	Forward(ctx context.Context, req *http.Request) (*http.Response, error)
	Host() string
}

type HttpServerProxy struct {
	target string
}

var defaultHttpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// Prevent redirects to ensure we handle them manually if needed
		return http.ErrUseLastResponse
	},
}

func (p *HttpServerProxy) Forward(ctx context.Context, req *http.Request) (*http.Response, error) {
	// rewrite the request to the target
	req.URL.Host = p.target
	req.URL.Scheme = "http"

	return defaultHttpClient.Do(req)
}

func (p *HttpServerProxy) Host() string {
	return p.target
}

func NewHttpServerProxy(target string) Proxy {
	return &HttpServerProxy{
		target: target,
	}
}
