package env

import "github.com/nitrictech/nitric/core/pkg/env"

var MONGODB_CONNECTION_STRING = env.GetEnv("MONGODB_CONNECTION_STRING", "")
var MONGODB_DATABASE = env.GetEnv("MONGODB_DATABASE", "")
var MONGODB_DIRECT = env.GetEnv("MONGODB_DIRECT", "true")

var KVAULT_NAME = env.GetEnv("KVAULT_NAME", "")

var AZURE_STORAGE_BLOB_ENDPOINT = env.GetEnv("AZURE_STORAGE_ACCOUNT_BLOB_ENDPOINT", "")
var AZURE_STORAGE_QUEUE_ENDPOINT = env.GetEnv("AZURE_STORAGE_ACCOUNT_QUEUE_ENDPOINT", "")

// mongoDBConnectionString := utils.GetEnv(mongoDBConnectionStringEnvVarName, "")

// 	if mongoDBConnectionString == "" {
// 		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBConnectionStringEnvVarName)
// 	}

// 	database := utils.GetEnv(mongoDBDatabaseEnvVarName, "")

// 	if database == "" {
// 		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBDatabaseEnvVarName)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	mongoDBSetDirect := utils.GetEnv(mongoDBSetDirectEnvVarName, "true")
