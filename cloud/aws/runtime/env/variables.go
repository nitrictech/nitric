package env

import "github.com/nitrictech/nitric/core/pkg/env"

// AWS_REGION - AWS Region where the application is currently executing.
var AWS_REGION = env.GetEnv("AWS_REGION", "us-east-1")

// GATEWAY_ENVIRONMENT - The environment the gateway needs to connect to e.g. "lambda" for aws lambda OR "http" for EC2
var GATEWAY_ENVIRONMENT = env.GetEnv("GATEWAY_ENVIRONMENT", "lambda")
