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

package sdk

import (
	"fmt"
	"github.com/valyala/fasthttp"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/triggers"
)

type NitricContext struct {
	RequestId   string
	PayloadType string
	Trigger     string
	TriggerType triggers.TriggerType
}

// Normalized NitricRequest
type NitricRequest struct {
	Context     *NitricContext
	ContentType string
	Payload     []byte
}

type NitricResponse struct {
	Headers map[string]string
	Status  int
	Body    []byte
}

type GatewayService interface {
	// Start the Gateway
	// This method should block
	Start(handler handler.TriggerHandler) (*fasthttp.Server, error)
}

type UnimplementedGatewayPlugin struct {
	GatewayService
}

func (*UnimplementedGatewayPlugin) Start(_ handler.TriggerHandler) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
