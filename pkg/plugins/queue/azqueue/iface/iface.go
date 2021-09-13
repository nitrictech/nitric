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

package azqueue_service_iface

import (
	"context"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type AzqueueServiceUrlIface interface {
	NewQueueURL(string) AzqueueQueueUrlIface
}

type AzqueueQueueUrlIface interface {
	NewMessageURL() AzqueueMessageUrlIface
}

type AzqueueMessageUrlIface interface {
	Enqueue(ctx context.Context, messageText string, visibilityTimeout time.Duration, timeToLive time.Duration) (*azqueue.EnqueueMessageResponse, error)
	Dequeue(ctx context.Context, maxMessages int32, visibilityTimeout time.Duration) (DequeueMessagesResponseIface, error)
	NewMessageIDURL(messageId azqueue.MessageID) AzqueueMessageIdUrlIface
}

type AzqueueMessageIdUrlIface interface {
	Delete(ctx context.Context, popReceipt azqueue.PopReceipt) (*azqueue.MessageIDDeleteResponse, error)
}

type DequeueMessagesResponseIface interface {
	NumMessages() int32
	Message(index int32) *azqueue.DequeuedMessage
}
