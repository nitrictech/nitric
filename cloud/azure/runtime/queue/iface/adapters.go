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

package iface

import (
	"context"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

func AdaptServiceUrl(c azqueue.ServiceURL) AzqueueServiceUrlIface {
	return serviceUrl{c}
}

func AdaptQueueUrl(c azqueue.QueueURL) AzqueueQueueUrlIface {
	return queueUrl{c}
}

func AdaptMessageUrl(c azqueue.MessagesURL) AzqueueMessageUrlIface {
	return messageUrl{c}
}

func AdaptMessageIdUrl(c azqueue.MessageIDURL) AzqueueMessageIdUrlIface {
	return messageIdUrl{c}
}

func AdaptDequeueMessagesResponse(c azqueue.DequeuedMessagesResponse) DequeueMessagesResponseIface {
	return dequeueMessagesResponse{c}
}

type (
	serviceUrl              struct{ c azqueue.ServiceURL }
	queueUrl                struct{ c azqueue.QueueURL }
	messageUrl              struct{ c azqueue.MessagesURL }
	messageIdUrl            struct{ c azqueue.MessageIDURL }
	dequeueMessagesResponse struct {
		c azqueue.DequeuedMessagesResponse
	}
)

func (c serviceUrl) NewQueueURL(queueName string) AzqueueQueueUrlIface {
	return AdaptQueueUrl(c.c.NewQueueURL(queueName))
}

func (c queueUrl) NewMessageURL() AzqueueMessageUrlIface {
	return AdaptMessageUrl(c.c.NewMessagesURL())
}

func (c messageUrl) Enqueue(ctx context.Context, messageText string, visibilityTimeout time.Duration, timeToLive time.Duration) (*azqueue.EnqueueMessageResponse, error) {
	return c.c.Enqueue(ctx, messageText, visibilityTimeout, timeToLive)
}

func (c messageUrl) Dequeue(ctx context.Context, maxMessages int32, visibilityTimeout time.Duration) (DequeueMessagesResponseIface, error) {
	resp, err := c.c.Dequeue(ctx, maxMessages, visibilityTimeout)
	if err != nil {
		return nil, err
	}
	return AdaptDequeueMessagesResponse(*resp), nil
}

func (c messageUrl) NewMessageIDURL(messageId azqueue.MessageID) AzqueueMessageIdUrlIface {
	return AdaptMessageIdUrl(c.c.NewMessageIDURL(messageId))
}

func (c messageIdUrl) Delete(ctx context.Context, popReceipt azqueue.PopReceipt) (*azqueue.MessageIDDeleteResponse, error) {
	return c.c.Delete(ctx, popReceipt)
}

func (c dequeueMessagesResponse) NumMessages() int32 {
	return c.c.NumMessages()
}

func (c dequeueMessagesResponse) Message(index int32) *azqueue.DequeuedMessage {
	return c.c.Message(index)
}
