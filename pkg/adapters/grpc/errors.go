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
	"fmt"

	v1 "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Provides GRPC error reporting
func NewGrpcError(operation string, err error) error {
	if pe, ok := err.(*errors.PluginError); ok {
		code := codes.Code(errors.Code(pe))

		ed := &v1.ErrorDetails{}
		ed.Message = pe.Msg
		if pe.Cause != nil {
			ed.Cause = pe.Cause.Error()
		}
		ed.Service = operation
		ed.Plugin = pe.Plugin
		ed.Args = pe.Args

		s := status.New(code, pe.Msg)
		s, _ = s.WithDetails(ed)

		return s.Err()

	} else {
		return newGrpcErrorWithCode(codes.Internal, operation, err)
	}
}

func newGrpcErrorWithCode(code codes.Code, operation string, err error) error {
	se := &v1.ErrorDetails{}
	se.Message = err.Error()
	se.Service = operation

	s := status.New(code, err.Error())
	s, _ = s.WithDetails(se)

	return s.Err()
}

// Provides generic error for unregistered plugins
func NewPluginNotRegisteredError(plugin string) error {
	se := &v1.ErrorDetails{}
	se.Message = fmt.Sprintf("%s plugin not registered", plugin)

	s := status.New(codes.Unimplemented, se.Message)
	s, _ = s.WithDetails(se)

	return s.Err()
}
