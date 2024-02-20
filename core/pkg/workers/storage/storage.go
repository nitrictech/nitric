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

package storage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nitrictech/nitric/core/pkg/help"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"github.com/nitrictech/nitric/core/pkg/workers"
)

// BucketName uniquely identifies a storage bucket
type BucketName = string

// BucketListenerManager manages storage listeners for different buckets
type BucketListenerManager struct {
	listenerMap map[BucketName][]*BucketEventListener
	mutex       sync.RWMutex
}

type BucketRequestHandler interface {
	storagepb.StorageListenerServer
	HandleRequest(request *storagepb.ServerMessage) (*storagepb.ClientMessage, error)
	WorkerCount() int
}

// WorkerConnection handles communication between storage and worker
type WorkerConnection = workers.WorkerRequestBroker[*storagepb.ServerMessage, *storagepb.ClientMessage]

// BucketEventListener listens for specific events on a storage bucket
type BucketEventListener struct {
	connection     *WorkerConnection
	bucketName     BucketName
	eventType      storagepb.BlobEventType
	keyPrefixMatch string
}

// WorkerCount returns the total number of workers across all listeners
func (b *BucketListenerManager) WorkerCount() int {
	total := 0
	for _, listeners := range b.listenerMap {
		total += len(listeners)
	}

	return total
}

// findMatchingListener or error if not found, for specific bucket, event type, and key prefix
func (b *BucketListenerManager) findMatchingListener(bucketName BucketName, eventType storagepb.BlobEventType, key string) (*BucketEventListener, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	listeners, exists := b.listenerMap[bucketName]
	if !exists {
		return nil, fmt.Errorf("no listeners registered for bucket %s", bucketName)
	}

	var matchedListener *BucketEventListener

	for _, listener := range listeners {
		if listener.eventType != eventType {
			continue
		}

		if strings.HasPrefix(key, listener.keyPrefixMatch) {
			matchedListener = listener
		}
	}

	if matchedListener == nil {
		return nil, fmt.Errorf("no listener registered for bucket %s and eventType %s with prefix matcher that matches blob key %s", bucketName, eventType, key)
	}

	return matchedListener, nil
}

// MutualPrefixCheck returns true if either string starts with the other
func MutualPrefixCheck(prefix, other string) bool {
	return strings.HasPrefix(prefix, other) || strings.HasPrefix(other, prefix)
}

var _ storagepb.StorageListenerServer = &BucketListenerManager{}

// RegisterNewListener adds a new listener for a given registration request
func (b *BucketListenerManager) RegisterNewListener(registration *storagepb.RegistrationRequest, stream storagepb.StorageListener_ListenServer) (*WorkerConnection, error) {
	bucketName := registration.GetBucketName()
	eventType := registration.GetBlobEventType()
	prefixFilter := registration.GetKeyPrefixFilter()
	if prefixFilter == "*" {
		prefixFilter = ""
	}

	workerConn := workers.NewWorkerRequestBroker[*storagepb.ServerMessage, *storagepb.ClientMessage](stream)

	if b.listenerMap[bucketName] == nil {
		b.listenerMap[bucketName] = make([]*BucketEventListener, 0, 1)
	} else {
		// Prevent overlapping prefixes for the same bucket and event type
		for _, existingListener := range b.listenerMap[bucketName] {
			if existingListener.eventType != eventType {
				continue
			}
			if MutualPrefixCheck(existingListener.keyPrefixMatch, prefixFilter) {
				return nil, fmt.Errorf("overlapping listener key prefixes %s and %s for bucket '%s'", existingListener.keyPrefixMatch, prefixFilter, bucketName)
			}
		}
	}

	newListener := &BucketEventListener{
		connection:     workerConn,
		bucketName:     bucketName,
		eventType:      eventType,
		keyPrefixMatch: prefixFilter,
	}
	b.listenerMap[bucketName] = append(b.listenerMap[bucketName], newListener)

	return workerConn, nil
}

// Listen handles incoming stream connections for storage events
func (b *BucketListenerManager) Listen(stream storagepb.StorageListener_ListenServer) error {
	initialMessage, err := stream.Recv()
	if err != nil {
		return err
	}

	registration := initialMessage.GetRegistrationRequest()
	if registration == nil {
		return fmt.Errorf("request received from unregistered storage listener, initial request must be a registration request. %s", help.BugInNitricHelpText())
	}

	workerConn, err := b.RegisterNewListener(registration, stream)
	if err != nil {
		return err
	}

	// Send acknowledgement of registration
	err = stream.Send(&storagepb.ServerMessage{
		Content: &storagepb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &storagepb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	return workerConn.Run()
}

// HandleRequest processes incoming requests and directs them to the appropriate listener
func (b *BucketListenerManager) HandleRequest(request *storagepb.ServerMessage) (*storagepb.ClientMessage, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	blobEventRequest := request.GetBlobEventRequest()
	if blobEventRequest == nil {
		return nil, fmt.Errorf("received unhandled request message type %T", request.Content)
	}

	bucketID := blobEventRequest.GetBucketName()
	eventType := blobEventRequest.GetBlobEvent().GetType()
	key := blobEventRequest.GetBlobEvent().GetKey()

	listener, err := b.findMatchingListener(bucketID, eventType, key)
	if err != nil {
		return nil, err
	}

	response, err := listener.connection.Send(request)

	return *response, err
}

func New() *BucketListenerManager {
	return &BucketListenerManager{
		listenerMap: make(map[BucketName][]*BucketEventListener),
		mutex:       sync.RWMutex{},
	}
}
