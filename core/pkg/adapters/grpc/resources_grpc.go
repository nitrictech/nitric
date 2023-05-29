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
	"fmt"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
)

type ResourcesServiceServer struct {
	v1.UnimplementedResourceServiceServer
	plugin resource.ResourceService
}

type ResourceServiceOption = func(*ResourcesServiceServer)

func WithResourcePlugin(plugin resource.ResourceService) ResourceServiceOption {
	return func(srv *ResourcesServiceServer) {
		if plugin != nil {
			srv.plugin = plugin
		}
	}
}

func (rs *ResourcesServiceServer) Declare(ctx context.Context, req *v1.ResourceDeclareRequest) (*v1.ResourceDeclareResponse, error) {
	err := rs.plugin.Declare(ctx, req)
	if err != nil {
		return nil, err
	}

	return &v1.ResourceDeclareResponse{}, nil
}

type resourceDetailsType interface {
	// Add union type here to exhaust possible ResourceDetailTypes
	*v1.ResourceDetailsResponse_Api
}

func convertDetails[T resourceDetailsType](details *resource.DetailsResponse[any]) (T, error) {
	switch det := details.Detail.(type) {
	case resource.ApiDetails:
		return &v1.ResourceDetailsResponse_Api{
			Api: &v1.ApiResourceDetails{
				Url: det.URL,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported details type")
	}
}

var resourceTypeMap = map[v1.ResourceType]resource.ResourceType{
	v1.ResourceType_Api: resource.ResourceType_Api,
}

func (rs *ResourcesServiceServer) Details(ctx context.Context, req *v1.ResourceDetailsRequest) (*v1.ResourceDetailsResponse, error) {
	cType, ok := resourceTypeMap[req.Resource.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", req.Resource.Type)
	}

	d, err := rs.plugin.Details(ctx, cType, req.Resource.Name)
	if err != nil {
		return nil, err
	}

	details, err := convertDetails(d)
	if err != nil {
		return nil, err
	}

	return &v1.ResourceDetailsResponse{
		Id:       d.Id,
		Provider: d.Provider,
		Service:  d.Service,
		Details:  details,
	}, nil
}

func NewResourcesServiceServer(opts ...ResourceServiceOption) v1.ResourceServiceServer {
	// Default server implementation
	srv := &ResourcesServiceServer{
		plugin: &resource.UnimplementResourceService{},
	}

	// Apply options
	for _, o := range opts {
		o(srv)
	}

	return srv
}
