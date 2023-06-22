package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// Represents a locally running http server
type HttpWorker struct {
	port int

	adapter.Adapter
}

var _ Worker = &HttpWorker{}

func (*HttpWorker) HandlesTrigger(req *v1.TriggerRequest) bool {
	// Can handle any given HTTP request
	return req.GetHttp() != nil
}

// TODO: We should proxy the request instead as best we can
// This will allow this worker to be added to a pool and used generically
// however we will want to make sure we don't miss anything in context
func (h *HttpWorker) HandleTrigger(ctx context.Context, req *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	targetHost, err := url.Parse(fmt.Sprintf("http://localhost:%d", h.port))
	if err != nil {
		return nil, err
	}

	newHeader := http.Header{}

	for k, v := range req.GetHttp().Headers {
		for _, val := range v.Value {
			if k == "X-Forwarded-Authorization" && newHeader["Authorization"] == nil {
				k = "Authorization"
			}
			newHeader.Add(k, val)
		}
	}

	targetPath := targetHost.JoinPath(req.GetHttp().Path)
	// forward the request to the server
	res, err := http.DefaultClient.Do(&http.Request{
		Header: newHeader,
		Method: req.GetHttp().GetMethod(),
		URL:    targetPath,
		Body:   io.NopCloser(bytes.NewReader(req.Data)),
	})
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseHeaders := map[string]*v1.HeaderValue{}

	for k, v := range res.Header {
		responseHeaders[k] = &v1.HeaderValue{
			Value: v,
		}
	}

	return &v1.TriggerResponse{
		Data: body,
		Context: &v1.TriggerResponse_Http{
			Http: &v1.HttpResponseContext{
				Status:  int32(res.StatusCode),
				Headers: responseHeaders,
			},
		},
	}, nil
}

func (h *HttpWorker) GetPort() int {
	return h.port
}

func NewHttpWorker(adapter adapter.Adapter, port int) *HttpWorker {
	return &HttpWorker{
		Adapter: adapter, 
		port: port,
	}
}
