package grpc

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
