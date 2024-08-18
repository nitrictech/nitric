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

package gateway

import (
	"fmt"

	apigateways "github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/jobs"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type GatewayStartOpts struct {
	ApiPlugin               apigateways.ApiRequestHandler
	HttpPlugin              http.HttpRequestHandler
	SchedulesPlugin         schedules.ScheduleRequestHandler
	TopicsListenerPlugin    topics.SubscriptionRequestHandler
	StorageListenerPlugin   storage.BucketRequestHandler
	WebsocketListenerPlugin websockets.WebsocketRequestHandler
	JobHandlerPlugin        jobs.JobRequestHandler
}

// GatewayService - The interface for a Nitric Gateway, which acts as provider specific adapter for all incoming requests.
type GatewayService interface {
	// Start the Gateway
	Start(opts *GatewayStartOpts) error
	// Stop the Gateway
	Stop() error
}

type UnimplementedGatewayPlugin struct {
	GatewayService
}

func (*UnimplementedGatewayPlugin) Start(_ *GatewayStartOpts) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedGatewayPlugin) Stop() error {
	return fmt.Errorf("UNIMPLEMENTED")
}
