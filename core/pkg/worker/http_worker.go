// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// A Nitric HTTP worker
type HttpWorker struct {
	address string
}

func (s *HttpWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	return true
}

func (h *HttpWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if http := trigger.GetHttp(); http != nil {
		address := fmt.Sprintf("http://%s%s", h.address, http.Path)

		httpRequest := fasthttp.AcquireRequest()
		httpRequest.SetRequestURI(address)

		for key, val := range http.QueryParams {
			for _, v := range val.Value {
				httpRequest.URI().QueryArgs().Add(key, v)
			}
		}

		for key, val := range http.Headers {
			for _, v := range val.Value {
				httpRequest.Header.Add(key, v)
			}
		}

		httpRequest.Header.Del("Content-Length")
		httpRequest.SetBody(trigger.Data)
		httpRequest.Header.SetContentLength(len(trigger.Data))

		var resp fasthttp.Response
		err := fasthttp.Do(httpRequest, &resp)
		if err != nil {
			return nil, err
		}

		headers := map[string]*v1.HeaderValue{}
		resp.Header.VisitAll(func(key []byte, val []byte) {
			headers[string(key)] = &v1.HeaderValue{
				Value: []string{string(val)},
			}
		})

		return &v1.TriggerResponse{
			Data: resp.Body(),
			Context: &v1.TriggerResponse_Http{
				Http: &v1.HttpResponseContext{
					Headers: headers,
					Status:  int32(resp.StatusCode()),
				},
			},
		}, nil
	} else if topic := trigger.GetTopic(); topic != nil {
		address := fmt.Sprintf("http://%s/subscriptions/%s", h.address, topic.Topic)

		httpRequest := fasthttp.AcquireRequest()
		httpRequest.SetRequestURI(address)

		var resp fasthttp.Response

		httpRequest.SetBody(trigger.Data)
		httpRequest.Header.SetContentLength(len(trigger.Data))

		err := fasthttp.Do(httpRequest, &resp)
		if err == nil && resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
			return &v1.TriggerResponse{
				Context: &v1.TriggerResponse_Topic{
					Topic: &v1.TopicResponseContext{
						Success: true,
					},
				},
			}, nil
		}
		if err != nil {
			return nil, errors.Wrapf(err, "Error processing event (%d): %s", resp.StatusCode(), string(resp.Body()))
		}
		return nil, errors.Errorf("Error processing event (%d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil, fmt.Errorf("invalid trigger provided")
}

// Creates a new HttpWorker
// Will wait to ensure that the provided address is dialable
// before proceeding
func NewHttpWorker(address string) (*HttpWorker, error) {
	// Dial the child port to see if it's open and ready...
	maxWaitTime := time.Duration(5) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	waitedTime := time.Duration(0)
	for {
		conn, _ := net.Dial("tcp", address)
		if conn != nil {
			conn.Close()
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return nil, fmt.Errorf("Unable to dial http worker, does it expose a http server at: %s?", address)
			}
		}
	}

	// Dial the provided address to ensure its availability
	return &HttpWorker{
		address: address,
	}, nil
}
