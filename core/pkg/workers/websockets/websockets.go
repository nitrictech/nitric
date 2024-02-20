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

package websockets

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nitrictech/nitric/core/pkg/help"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	workers "github.com/nitrictech/nitric/core/pkg/workers"
)

// WorkerConnection manages communication between websocket and worker
type WorkerConnection = workers.WorkerRequestBroker[*websocketspb.ServerMessage, *websocketspb.ClientMessage]

type WebsocketRequestHandler interface {
	websocketspb.WebsocketHandlerServer
	HandleRequest(request *websocketspb.ServerMessage) (*websocketspb.ClientMessage, error)
	WorkerCount() int
}

// WebsocketManager manages connections and event handlers for websockets
type WebsocketManager struct {
	handlers map[string]*WorkerConnection
	mutex    sync.RWMutex
}

// generateHandlerKey creates a unique identifier for a websocket event handler
func generateHandlerKey(socketName string, eventType websocketspb.WebsocketEventType) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", socketName, eventType.String()))
}

// WorkerCount returns the total number of websocket handlers
func (wm *WebsocketManager) WorkerCount() int {
	return len(wm.handlers)
}

// registerHandler adds a new handler to the manager
func (wm *WebsocketManager) registerHandler(handler *WorkerConnection, registrationRequest *websocketspb.RegistrationRequest) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	socketName := registrationRequest.SocketName
	eventType := registrationRequest.EventType

	handlerKey := generateHandlerKey(socketName, eventType)

	if _, exists := wm.handlers[handlerKey]; exists {
		return fmt.Errorf("websocket handler already registered, socket: %s eventType: %s", socketName, eventType.String())
	}

	wm.handlers[handlerKey] = handler
	return nil
}

func (wm *WebsocketManager) unregisterHandler(handler *WorkerConnection) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	var resultKey string

	for k, wrb := range wm.handlers {
		if wrb == handler {
			resultKey = k
			break
		}
	}

	delete(wm.handlers, resultKey)
}

// ManageEventHandlers handles the registration of new websocket event handlers
func (wm *WebsocketManager) HandleEvents(stream websocketspb.WebsocketHandler_HandleEventsServer) error {
	initialMessage, err := stream.Recv()
	if err != nil {
		return err
	}

	registration := initialMessage.GetRegistrationRequest()
	if registration == nil {
		return fmt.Errorf("unregistered websocket handler, expected a registration request. %s", help.BugInNitricHelpText())
	}

	handler := workers.NewWorkerRequestBroker[*websocketspb.ServerMessage, *websocketspb.ClientMessage](stream)
	if err := wm.registerHandler(handler, registration); err != nil {
		return err
	}

	defer wm.unregisterHandler(handler)

	err = stream.Send(&websocketspb.ServerMessage{
		Content: &websocketspb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &websocketspb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	return handler.Run()
}

// FindMatchingHandler returns a handler for a specific socket and event type, or an error if not found
func (wm *WebsocketManager) FindMatchingHandler(socketName string, eventType websocketspb.WebsocketEventType) (*WorkerConnection, error) {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	handlerKey := generateHandlerKey(socketName, eventType)

	handler, exists := wm.handlers[handlerKey]
	if !exists {
		return nil, fmt.Errorf("no handlers for socket: %s and eventType: %s", socketName, eventType.String())
	}

	return handler, nil
}

// HandleRequest handles incoming requests and forwards them to the appropriate handler
func (wm *WebsocketManager) HandleRequest(request *websocketspb.ServerMessage) (*websocketspb.ClientMessage, error) {
	eventRequest := request.GetWebsocketEventRequest()
	if eventRequest == nil {
		return nil, fmt.Errorf("invalid request, expected a websocket event request. %s", help.BugInNitricHelpText())
	}

	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	socketName := eventRequest.SocketName
	eventType := determineEventType(eventRequest)

	handler, err := wm.FindMatchingHandler(socketName, eventType)
	if err != nil {
		return nil, err
	}

	response, err := handler.Send(request)
	if err != nil {
		return nil, err
	}

	return *response, nil
}

// determineEventType deduces the event type from the event request
func determineEventType(eventRequest *websocketspb.WebsocketEventRequest) websocketspb.WebsocketEventType {
	if eventRequest.GetDisconnection() != nil {
		return websocketspb.WebsocketEventType_Disconnect
	} else if eventRequest.GetMessage() != nil {
		return websocketspb.WebsocketEventType_Message
	}
	return websocketspb.WebsocketEventType_Connect
}

// NewWebsocketManager creates a new instance of WebsocketManager
func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		handlers: make(map[string]*WorkerConnection),
		mutex:    sync.RWMutex{},
	}
}
