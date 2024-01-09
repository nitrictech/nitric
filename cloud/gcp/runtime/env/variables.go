package env

import "github.com/nitrictech/nitric/core/pkg/env"

var GCP_REGION = env.GetEnv("GCP_REGION", "")
var GOOGLE_PROJECT_ID = env.GetEnv("GOOGLE_PROJECT_ID", "")
var SERVICE_ACCOUNT_EMAIL = env.GetEnv("SERVICE_ACCOUNT_EMAIL", "")

// The name of the google cloud tasks queue to use to delay message delivery to pubsub topics
var DELAY_QUEUE_NAME = env.GetEnv("DELAY_QUEUE_NAME", "")
