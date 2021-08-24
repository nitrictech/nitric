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
	Dequeue(ctx context.Context, maxMessages int32, visibilityTimeout time.Duration) (*azqueue.DequeuedMessagesResponse, error)
	NewMessageIDURL(messageId azqueue.MessageID) AzqueueMessageIdUrlIface
}

type AzqueueMessageIdUrlIface interface {
	Delete(ctx context.Context, popReceipt azqueue.PopReceipt) (*azqueue.MessageIDDeleteResponse, error)
}
