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

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// WebsockerWorker - Worker representation for a websocket handler
type WebsocketWorker struct {
	socket string
	event  v1.WebsocketEvent

	adapter.Adapter
}

var _ Worker = &RouteWorker{}

func (s *WebsocketWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	if websocket := trigger.GetWebsocket(); websocket != nil {
		return websocket.Socket == s.socket && websocket.Event == s.event
	}

	return false
}

// Socket - Retrieve the name of the socket this
// websocket worker was registered for
func (s *WebsocketWorker) Socket() string {
	return s.socket
}

func (s *WebsocketWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if http := trigger.GetWebsocket(); http != nil {
		return s.Adapter.HandleTrigger(ctx, trigger)
	}

	return nil, fmt.Errorf("websocket worker does not handle non-Websocket triggers")
}

type WebsocketWorkerOptions struct {
	Socket string
	Event  v1.WebsocketEvent
}

func NewWebsocketWorker(adapter adapter.Adapter, opts *WebsocketWorkerOptions) *WebsocketWorker {
	return &WebsocketWorker{
		socket:  opts.Socket,
		event:   opts.Event,
		Adapter: adapter,
	}
}
