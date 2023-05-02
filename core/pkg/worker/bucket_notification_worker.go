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
	"regexp"
	"sort"
	"strings"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// BucketNotificationWorker - Worker representation for an bucket notification handler
type BucketNotificationWorker struct {
	notification *v1.BucketNotificationWorker
	adapter.Adapter
}

var _ Worker = &BucketNotificationWorker{}

func (s *BucketNotificationWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	notification := trigger.GetNotification()
	if notification == nil {
		return false
	}

	if s.notification.Bucket != notification.Source {
		return false
	}

	if !s.matchesPrefixFilter(notification.GetBucket().Key) {
		return false
	}

	if s.notification.Config.NotificationType != notification.GetBucket().Type {
		return false
	}

	return true
}

func (s *BucketNotificationWorker) Bucket() string {
	return s.notification.Bucket
}

func (s *BucketNotificationWorker) NotificationType() v1.BucketNotificationType {
	return s.notification.Config.NotificationType
}

func (s *BucketNotificationWorker) NotificationPrefixFilter() string {
	return s.notification.Config.NotificationPrefixFilter
}

func (s *BucketNotificationWorker) matchesPrefixFilter(objectKey string) bool {
	eventFilter := s.NotificationPrefixFilter()
	if eventFilter == "*" {
		eventFilter = ""
	}

	eventFilterRegex := fmt.Sprintf("^(%s)", strings.ReplaceAll(eventFilter, "/", "\\/"))
	match, _ := regexp.MatchString(eventFilterRegex, objectKey)

	return match
}

func (s *BucketNotificationWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if trigger.GetNotification() == nil {
		return nil, fmt.Errorf("cannot handle given notification")
	}

	return s.Adapter.HandleTrigger(ctx, trigger)
}

type BucketNotificationWorkerOptions struct {
	Notification *v1.BucketNotificationWorker
}

// Checks that there are no overlapping bucket notifications
func ValidateBucketNotifications(workers []Worker) error {
	notificationByEventType := make(map[string]map[v1.BucketNotificationType][]string)

	// Filter for only notification workers
	notifications := []*v1.BucketNotificationWorker{}
	for _, w := range workers {
		if notificationWrkr, ok := w.(*BucketNotificationWorker); ok {
			notifications = append(notifications, notificationWrkr.notification)
		}
	}

	// Separate the notifications by event type and bucket name
	for _, n := range notifications {
		eventFilter := n.Config.NotificationPrefixFilter
		if eventFilter == "*" {
			eventFilter = ""
		}

		if notificationByEventType[n.Bucket] == nil {
			notificationByEventType[n.Bucket] = make(map[v1.BucketNotificationType][]string)
		}

		notificationByEventType[n.Bucket][n.Config.NotificationType] = append(notificationByEventType[n.Bucket][n.Config.NotificationType], eventFilter)
	}

	for bucketName := range notificationByEventType {
		for _, eventType := range []v1.BucketNotificationType{v1.BucketNotificationType_Created, v1.BucketNotificationType_Deleted} {
			// Sort by the path
			events := notificationByEventType[bucketName][eventType]
			sort.Strings(events)

			for idx, n := range events {
				if n == events[len(events)-1] {
					break
				}

				match, err := regexp.MatchString(fmt.Sprintf("^(%s)", strings.ReplaceAll(n, "/", "\\/")), events[idx+1])
				if err != nil {
					return err
				}

				if match {
					return fmt.Errorf("overlapping prefixes in notifications for bucket '%s'", bucketName)
				}
			}
		}
	}

	return nil
}

// Package private method
// Only a pool may create a new faas worker
func NewBucketNotificationWorker(adapter adapter.Adapter, opts *BucketNotificationWorkerOptions) *BucketNotificationWorker {
	return &BucketNotificationWorker{
		notification: opts.Notification,
		Adapter:      adapter,
	}
}
