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

// AWS_REGION - AWS Region where the application is currently executing.
var AWS_REGION = env.GetEnv("AWS_REGION", "us-east-1")

// GATEWAY_ENVIRONMENT - The environment the gateway needs to connect to e.g. "lambda" for aws lambda OR "http" for EC2
var GATEWAY_ENVIRONMENT = env.GetEnv("GATEWAY_ENVIRONMENT", "lambda")

// JOB_QUEUE_ARN - The AWS ARN of the job queue to use for job execution
var JOB_QUEUE_ARN = env.GetEnv("NITRIC_JOB_QUEUE_ARN", "")

var NITRIC_AWS_RESOURCE_RESOLVER = env.GetEnv("NITRIC_AWS_RESOURCE_RESOLVER", "ssm")
