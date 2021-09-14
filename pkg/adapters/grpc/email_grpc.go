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

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/pkg/plugins/emails"
	"google.golang.org/grpc/codes"
)

// EmailServiceServer - GRPC Interface for registered Nitric Email Plugins
type EmailServiceServer struct {
	pb.UnimplementedEmailServiceServer
	emailPlugin emails.EmailService
}

func (s *EmailServiceServer) checkPluginRegistered() error {
	if s.emailPlugin == nil {
		return NewPluginNotRegisteredError("Email")
	}

	return nil
}

// Send - Send an email using the provided inputs
func (s *EmailServiceServer) Send(ctx context.Context, req *pb.EmailSendRequest) (*pb.EmailSendResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "EmailService.Send", err)
	}

	dest := emails.EmailDestination{
		To:  req.To,
		Cc:  req.Cc,
		Bcc: req.Bcc,
	}

	body := emails.EmailBody{
		Text: &req.Body.Text,
		Html: &req.Body.Html,
	}

	if err := s.emailPlugin.Send(req.From, dest, req.Subject, body); err == nil {
		return &pb.EmailSendResponse{}, nil
	} else {
		return nil, NewGrpcError("EmailService.Send", err)
	}
}

func NewEmailServiceServer(emailPlugin emails.EmailService) pb.EmailServiceServer {
	return &EmailServiceServer{
		emailPlugin: emailPlugin,
	}
}
