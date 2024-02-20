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
	MONGODB_CONNECTION_STRING = env.GetEnv("MONGODB_CONNECTION_STRING", "")
	MONGODB_DATABASE          = env.GetEnv("MONGODB_DATABASE", "")
	MONGODB_DIRECT            = env.GetEnv("MONGODB_DIRECT", "true")
)

var KVAULT_NAME = env.GetEnv("KVAULT_NAME", "")

var AZURE_STORAGE_ACCOUNT_NAME = env.GetEnv("AZURE_STORAGE_ACCOUNT_NAME", "")

var (
	AZURE_STORAGE_BLOB_ENDPOINT  = env.GetEnv("AZURE_STORAGE_ACCOUNT_BLOB_ENDPOINT", "")
	AZURE_STORAGE_QUEUE_ENDPOINT = env.GetEnv("AZURE_STORAGE_ACCOUNT_QUEUE_ENDPOINT", "")
)

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
