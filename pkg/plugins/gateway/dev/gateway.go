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

// The AWS HTTP gateway plugin
package gateway_plugin

import (
	"fmt"
	"strings"

	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/worker"

	"github.com/nitric-dev/membrane/pkg/plugins/gateway"
	"github.com/nitric-dev/membrane/pkg/plugins/gateway/base_http"
	"github.com/valyala/fasthttp"
)

func middleware(ctx *fasthttp.RequestCtx, wrkr worker.Worker) bool {
	var triggerTypeString = string(ctx.Request.Header.Peek("x-nitric-source-type"))

	// Handle Event/Subscription Request Types
	if strings.ToUpper(triggerTypeString) == triggers.TriggerType_Subscription.String() {
		trigger := string(ctx.Request.Header.Peek("x-nitric-source"))
		requestId := string(ctx.Request.Header.Peek("x-nitric-request-id"))
		payload := ctx.Request.Body()

		err := wrkr.HandleEvent(&triggers.Event{
			ID:      requestId,
			Topic:   trigger,
			Payload: payload,
		})

		if err != nil {
			fmt.Println(err)
			ctx.Error(fmt.Sprintf("Error processing event. Details: %s", err), 500)
		} else {
			ctx.SuccessString("text/plain", "Successfully Handled the Event")
		}

		// return here...
		return false
	}

	return true
}

// Create new HTTP gateway
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (gateway.GatewayService, error) {
	return base_http.New(middleware)
}
