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

type Delegate interface {
	HandlesTrigger(*v1.TriggerRequest) bool
}

type Worker interface {
	Delegate
	adapter.Adapter
}

type UnimplementedWorker struct{}

func (*UnimplementedWorker) HandlesTrigger(trig *v1.TriggerRequest) bool {
	return false
}

func (*UnimplementedWorker) HandleTrigger(ctx context.Context, trig *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	return nil, fmt.Errorf("unimplemented worker does not handle triggers")
}
