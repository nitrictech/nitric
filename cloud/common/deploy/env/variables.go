package env

import "github.com/nitrictech/nitric/core/pkg/env"

var PORT = env.GetEnv("PORT", "50051")
