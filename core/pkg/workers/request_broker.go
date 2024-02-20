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

package workers

import (
	"fmt"
	"sync"

	"github.com/nitrictech/nitric/core/pkg/help"
	"github.com/nitrictech/nitric/core/pkg/logger"
	"google.golang.org/grpc"
)

type GrpcBidiStreamServer[ServerMessage IdentifiableMessage, ClientMessage IdentifiableMessage] interface {
	Send(ServerMessage) error
	Recv() (ClientMessage, error)
	grpc.ServerStream
}

type RequestIdentifier = string

// IdentifiableMessage is a message that has an ID, used to match requests and responses
type IdentifiableMessage interface {
	GetId() RequestIdentifier
}

// WorkerRequestBroker helps manage the async bidirectional stream between the worker (typically a client SDK) and the Nitric server.
//
//	the broker facilitates sending requests to a worker and awaits responses, then matches them with the corresponding request.
//	This enables users of this brokers to treat the request/response lifecycle with the worker as if they were synchronous.
type WorkerRequestBroker[Request IdentifiableMessage, Response IdentifiableMessage] struct {
	workerConnectionStream GrpcBidiStreamServer[Request, Response]
	responseChannelLock    sync.RWMutex
	responseChannels       map[RequestIdentifier]chan Response
	running                bool
}

func (w *WorkerRequestBroker[Request, Response]) Send(req Request) (*Response, error) {
	if !w.running {
		return nil, fmt.Errorf("worker server not running, call Start() before sending requests")
	}

	if _, exists := w.responseChannels[req.GetId()]; exists {
		return nil, fmt.Errorf("request with ID %s already exists", req.GetId())
	}

	w.responseChannelLock.Lock()
	responseChannel := make(chan Response)
	w.responseChannels[req.GetId()] = responseChannel
	w.responseChannelLock.Unlock()

	if err := w.workerConnectionStream.Send(req); err != nil {
		return nil, err
	}

	// wait for the response
	response, ok := <-responseChannel
	if !ok {
		return nil, fmt.Errorf("error receiving response")
	}

	// clean up the map reference
	w.responseChannelLock.Lock()
	delete(w.responseChannels, req.GetId())
	w.responseChannelLock.Unlock()

	return &response, nil
}

// Run the connection broker, allowing async communication with the worker (i.e.
func (w *WorkerRequestBroker[Request, Response]) Run() error {
	if w.running {
		return fmt.Errorf("worker already running")
	}
	// Read responses on the client connection stream and match them with the corresponding response channel.
	for {
		w.running = true
		response, err := w.workerConnectionStream.Recv()
		if err != nil {
			// Most likely the client closed the connection
			return err
		}

		w.responseChannelLock.RLock()
		responseChannel, ok := w.responseChannels[response.GetId()]
		w.responseChannelLock.RUnlock()
		if !ok {
			// This would indicate a critical bug, it means that the client (SDK) did not return a response with an ID that matches a request that was sent to it
			// OR there may have been a network error that resulted in duplicate responses
			logger.Errorf("nitric received a response for an unknown request, response could not be returned: %s", help.BugInNitricHelpText())
		}

		responseChannel <- response
	}
}

func NewWorkerRequestBroker[Request IdentifiableMessage, Response IdentifiableMessage](workerConnectionStream GrpcBidiStreamServer[Request, Response]) *WorkerRequestBroker[Request, Response] {
	return &WorkerRequestBroker[Request, Response]{
		workerConnectionStream: workerConnectionStream,
		responseChannelLock:    sync.RWMutex{},
		responseChannels:       make(map[string]chan Response),
		running:                false,
	}
}
