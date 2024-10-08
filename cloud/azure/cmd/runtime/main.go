// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/nitrictech/nitric/cloud/azure/runtime"
	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	"github.com/nitrictech/nitric/core/pkg/logger"
	"github.com/nitrictech/nitric/core/pkg/server"
)

func main() {
	resourceResolver, err := resource.New()
	if err != nil {
		logger.Fatalf("could not create core azure resource resolver: %v", err)
	}

	m, err := runtime.NewAzureRuntimeServer(resourceResolver)
	if err != nil {
		logger.Fatalf("There was an error initializing the nitric server: %v", err)
	}

	server.Run(m)
}
