package service

import (
	"context"
	"net/http"
)

type Proxy interface {
	Forward(ctx context.Context, req *http.Request) (*http.Response, error)
}

type HttpServerProxy struct {
	target string
}

func (p *HttpServerProxy) Forward(ctx context.Context, req *http.Request) (*http.Response, error) {
	// rewrite the request to the target
	req.URL.Host = p.target
	req.URL.Scheme = "http"

	return http.DefaultClient.Do(req)
}

func NewHttpServerProxy(target string) Proxy {
	return &HttpServerProxy{
		target: target,
	}
}
