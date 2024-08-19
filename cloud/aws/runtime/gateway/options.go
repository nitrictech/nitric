// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gateway

type lambdaGatewayOption func(*LambdaGateway)

// WithRuntime sets the lambda runtime handler for the LambdaGateway
func WithRuntime(runtime LambdaRuntimeHandler) lambdaGatewayOption {
	return func(g *LambdaGateway) {
		g.runtime = runtime
	}
}

// WithRouter sets the lambda event router for the LambdaGateway
func WithRouter(router LambdaEventRouter) lambdaGatewayOption {
	return func(g *LambdaGateway) {
		g.routeEvent = router
	}
}
