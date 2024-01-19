package env

import "github.com/nitrictech/nitric/core/pkg/env"

// The unique id of the nitric stack that this service in running in
var NITRIC_STACK_ID = env.GetEnv("NITRIC_STACK_ID", "")

// % of	calls to trace, 0-100
var NITRIC_TRACE_SAMPLE_PERCENT = env.GetEnv("NITRIC_TRACE_SAMPLE_PERCENT", "0")

// Address of the Gateway to register for
var GATEWAY_ADDRESS = env.GetEnv("GATEWAY_ADDRESS", ":9001")
