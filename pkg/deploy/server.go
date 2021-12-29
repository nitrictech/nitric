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

package deploy

import (
	"context"
	"fmt"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	function *Function
	pb.UnimplementedFaasServiceServer
	pb.UnimplementedResourceServiceServer
}

// TriggerStream - Starts a new FaaS server stream
//
// The deployment server collects information from stream InitRequests, then immediately terminates the stream
// This behavior captures enough information to identify function handlers, without executing the handler code
// during the build process.
func (s *Server) TriggerStream(stream pb.FaasService_TriggerStreamServer) error {
	cm, err := stream.Recv()

	if err != nil {
		return status.Errorf(codes.Internal, "error reading message from stream: %v", err)
	}

	ir := cm.GetInitRequest()

	if ir == nil {
		// SHUT IT DOWN!!!!
		// The first message must be an init request from the prospective FaaS worker
		return status.Error(codes.FailedPrecondition, "first message must be InitRequest")
	}

	switch w := ir.Worker.(type) {
	case *pb.InitRequest_Api:
		s.function.AddApiHandler(w.Api)
	case *pb.InitRequest_Schedule:
		s.function.AddScheduleHandler(w.Schedule)
	case *pb.InitRequest_Subscription:
		s.function.AddSubscriptionHandler(w.Subscription)
	default:
		// treat as normal function worker
		// XXX: No-op for now. This can be handled exclusively at runtime
	}

	fmt.Println(s.function.String())

	// Close the stream, once we've recieved the InitRequest
	return nil
}

// Declare - Accepts resource declarations, adding them as dependencies to the Function
func (s *Server) Declare(ctx context.Context, req *pb.ResourceDeclareRequest) (*pb.ResourceDeclareResponse, error) {

	switch req.Resource.Type {
	case pb.ResourceType_Bucket:
		s.function.AddBucket(req.Resource.Name, req.GetBucket())
	case pb.ResourceType_Collection:
		s.function.AddCollection(req.Resource.Name, req.GetCollection())
	case pb.ResourceType_Queue:
		s.function.AddQueue(req.Resource.Name, req.GetQueue())
	case pb.ResourceType_Topic:
		s.function.AddTopic(req.Resource.Name, req.GetTopic())
	}

	fmt.Println(s.function.String())

	tmpStack := &Stack{functions: []*Function{s.function}}
	for a, _ := range s.function.apis {
		if spec, e := tmpStack.GetApiSpec(a); spec != nil {
			fmt.Println("oaiSpec", spec)
		} else {
			fmt.Println("specErr", e.Error())
		}
	}

	return &pb.ResourceDeclareResponse{}, nil
}

// New - Creates a new deployment server
func New(function *Function) *Server {
	return &Server{
		function: function,
	}
}