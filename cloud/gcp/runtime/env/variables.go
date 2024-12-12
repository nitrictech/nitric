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

var (
	GCP_REGION              = env.GetEnv("GCP_REGION", "")
	GOOGLE_PROJECT_ID       = env.GetEnv("GOOGLE_PROJECT_ID", "")
	SERVICE_ACCOUNT_EMAIL   = env.GetEnv("SERVICE_ACCOUNT_EMAIL", "")
	JOBS_BUCKET_NAME        = env.GetEnv("NITRIC_JOBS_BUCKET_NAME", "")
	FIRESTORE_DATABASE_NAME = env.GetEnv("FIRESTORE_DATABASE_NAME", "(default)")
)

// The name of the google cloud tasks queue to use to delay message delivery to pubsub topics
var DELAY_QUEUE_NAME = env.GetEnv("DELAY_QUEUE_NAME", "")
