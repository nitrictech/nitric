// Copyright 2021 Nitric Technologies Pty Ltd.
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
	"strings"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// RouteWorker - Worker representation for an http api route handler
type BucketNotificationWorker struct {
	bucket string
	adapter.Adapter
}

var _ Worker = &BucketNotificationWorker{}

func NotificationBucketToBucketName(bucket string) string {
	return strings.ToLower(strings.ReplaceAll(bucket, " ", "-"))
}

func (s *BucketNotificationWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	if notification := trigger.GetNotification(); notification != nil {
		return notification.Type == v1.NotificationType_Bucket && NotificationBucketToBucketName(s.bucket) == notification.Resource
	}

	return false
}

func (s *BucketNotificationWorker) Bucket() string {
	return s.bucket
}

func (s *BucketNotificationWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if trigger.GetNotification() == nil {
		return nil, fmt.Errorf("cannot handle given notification")
	}

	return s.Adapter.HandleTrigger(ctx, trigger)
}

type BucketNotificationWorkerOptions struct {
	Bucket string
}

// Package private method
// Only a pool may create a new faas worker
func NewBucketNotificationWorker(adapter adapter.Adapter, opts *BucketNotificationWorkerOptions) *BucketNotificationWorker {
	return &BucketNotificationWorker{
		bucket:  opts.Bucket,
		Adapter: adapter,
	}
}
