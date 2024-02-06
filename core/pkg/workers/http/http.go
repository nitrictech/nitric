// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"net"
	"sync"
	"time"

	httppb "github.com/nitrictech/nitric/core/pkg/proto/http/v1"
	"github.com/valyala/fasthttp"
)

type HttpServer struct {
	Output chan error
	host   string
	lock   sync.Mutex
}

type HttpRequestHandler interface {
	httppb.HttpServer
	HandleRequest(request *fasthttp.Request) (*fasthttp.Response, error)
	WorkerCount() int
}

// IsPortOpen returns true if the port is open, false otherwise
// returns false if the timeout elapses before a connection is established
func IsPortOpen(host string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false
	}

	err = conn.Close()
	return err == nil
}

func (h *HttpServer) WorkerCount() int {
	if h.host != "" {
		return 1
	}

	return 0
}

const (
	// Dial timeout for initial server check
	initialStartTimeout = 5 * time.Second
	// Polling interval for liveness check
	// portPollInterval = 5 * time.Second
	// // Dial timeout when polling
	// subsequentTimeout = 25 * time.Millisecond
)

func (h *HttpServer) Proxy(stream httppb.Http_ProxyServer) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.host != "" {
		return fmt.Errorf("http server already registered")
	}

	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	if msg.GetRequest() == nil {
		return fmt.Errorf("first request must be a proxy request")
	}

	host := msg.Request.Host

	if host == "" {
		return fmt.Errorf("host is required")
	}

	if !IsPortOpen(host, initialStartTimeout) {
		return fmt.Errorf("host %s failed to respond within %s", host, initialStartTimeout.String())
	}

	h.host = host
	for {
		_, err = stream.Recv()

		if err != nil {
			break
		}
	}

	h.host = ""
	h.Output <- fmt.Errorf("HTTP server %s is not available", host)
	return err
}

// HandleRequest forwards proxy request to the underlying HTTP server
func (srv *HttpServer) HandleRequest(request *fasthttp.Request) (*fasthttp.Response, error) {
	if srv.host == "" {
		return nil, fmt.Errorf("http server not registered")
	}

	requestCopy := &fasthttp.Request{}
	var response fasthttp.Response

	request.CopyTo(requestCopy)
	requestCopy.URI().SetHost(srv.host)
	requestCopy.URI().SetScheme("http")

	if err := fasthttp.Do(requestCopy, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func New() *HttpServer {
	return &HttpServer{
		Output: make(chan error),
	}
}
