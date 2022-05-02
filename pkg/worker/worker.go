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
	"fmt"

	"github.com/nitrictech/nitric/pkg/triggers"
)

type Worker interface {
	HandleEvent(trigger *triggers.Event) error
	HandleCloudEvent(trigger *triggers.CloudEvent) error
	HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error)
	HandlesHttpRequest(trigger *triggers.HttpRequest) bool
	HandlesEvent(trigger *triggers.Event) bool
	HandlesCloudEvent(trigger *triggers.CloudEvent) bool
}

type UnimplementedWorker struct{}

func (*UnimplementedWorker) HandlesEvent(trigger *triggers.Event) bool {
	return false
}

func (*UnimplementedWorker) HandlesCloudEvent(trigger *triggers.CloudEvent) bool {
	return false
}

func (*UnimplementedWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return false
}

func (*UnimplementedWorker) HandleEvent(trigger *triggers.Event) error {
	return fmt.Errorf("worker does not handle events")
}

func (*UnimplementedWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	return nil, fmt.Errorf("worker does not handle http requests")
}

func (*UnimplementedWorker) HandleCloudEvent(trigger *triggers.CloudEvent) (*triggers.CloudEvent, error) {
	return nil, fmt.Errorf("worker does not handle cloud events")
}
