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

package env

import "github.com/nitrictech/nitric/core/pkg/env"

// The unique id of the nitric stack that this service in running in
var NITRIC_STACK_ID = env.GetEnv("NITRIC_STACK_ID", "")

// % of	calls to trace, 0-100
var NITRIC_TRACE_SAMPLE_PERCENT = env.GetEnv("NITRIC_TRACE_SAMPLE_PERCENT", "0")

// Address of the Gateway to register for
var GATEWAY_ADDRESS = env.GetEnv("GATEWAY_ADDRESS", ":9001")
