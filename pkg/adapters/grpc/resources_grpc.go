package grpc

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

import (
	"context"

	v1 "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

type ResourcesServiceServer struct {
	v1.UnimplementedResourceServiceServer
}

func (rs *ResourcesServiceServer) Declare(ctx context.Context, req *v1.ResourceDeclareRequest) (*v1.ResourceDeclareResponse, error) {
	// Currently a no-op at runtime
	// TODO: Implement a strategy pattern for resolving resources, by their declared resource name in nitric
	return &v1.ResourceDeclareResponse{}, nil
}

func NewResourcesServiceServer() v1.ResourceServiceServer {
	return &ResourcesServiceServer{}
}
