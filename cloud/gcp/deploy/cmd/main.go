// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nitrictech/nitric/cloud/gcp/deploy"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"google.golang.org/grpc"
)

// Start the deployment server
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error listening on port 50051 %v", err)
	}

	srv := grpc.NewServer()

	deploySrv, err := deploy.NewServer()
	if err != nil {
		log.Fatalf("error starting deployment server %v", err)
	}
	v1.RegisterDeployServiceServer(srv, deploySrv)

	fmt.Printf("Deployment server started on %s\n", "") // lis.Addr().String()
	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("error serving requests %v", err)
	}
}
