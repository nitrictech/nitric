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

package grpc

import (
	"context"

	"google.golang.org/grpc/codes"

	pb "github.com/nitrictech/nitric/core/pkg/api/nitric/websocket/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/websocket"
)

// GRPC Interface for registered Nitric Storage Plugins
type WebsocketServiceServer struct {
	pb.UnimplementedWebsocketServiceServer
	websocketPlugin websocket.WebsocketService
}

func (s *WebsocketServiceServer) checkPluginRegistered() error {
	if s.websocketPlugin == nil {
		return NewPluginNotRegisteredError("Websocket")
	}

	return nil
}

func (s *WebsocketServiceServer) Send(ctx context.Context, req *pb.WebsocketSendRequest) (*pb.WebsocketSendResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "WebsocketService.Send", err)
	}

	if err := s.websocketPlugin.Send(ctx, req.Socket, req.ConnectionId, req.Data); err == nil {
		return &pb.WebsocketSendResponse{}, nil
	} else {
		return nil, NewGrpcError("WebsocketService.Send", err)
	}
}

func (s *WebsocketServiceServer) Close(ctx context.Context, req *pb.WebsocketCloseRequest) (*pb.WebsocketCloseResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "WebsocketService.Close", err)
	}

	if err := s.websocketPlugin.Close(ctx, req.Socket, req.ConnectionId); err == nil {
		return &pb.WebsocketCloseResponse{}, nil
	} else {
		return nil, NewGrpcError("WebsocketService.Close", err)
	}
}

func NewWebsocketServiceServer(websocketPlugin websocket.WebsocketService) pb.WebsocketServiceServer {
	return &WebsocketServiceServer{
		websocketPlugin: websocketPlugin,
	}
}
