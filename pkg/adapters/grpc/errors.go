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
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcError struct {
	code    codes.Code `json:"-"`
	Code    string     `json:"code"`
	Msg     string     `json:"msg,omitempty"`
	Service string     `json:"service,omitempty"`
}

func (ge *grpcError) Error() string {
	ge.Code = fmt.Sprintf("%v", ge.code)
	data, _ := json.Marshal(ge)
	return string(data)
}

// Provides GRPC error reporting
func NewGrpcError(operation string, err error) error {
	if pe, ok := err.(*errors.PluginError); ok {
		pe.Service = operation
		return status.Error(codes.Code(errors.Code(pe)), pe.Error())
	} else {
		return newGrpcErrorWithCode(codes.Internal, operation, err)
	}
}

func newGrpcErrorWithCode(code codes.Code, operation string, err error) error {
	ge := &grpcError{
		Service: operation,
		code:    code,
		Msg:     err.Error(),
	}
	return status.Error(code, ge.Error())
}

// Provides generic error for unregistered plugins
func NewPluginNotRegisteredError(plugin string) error {
	ge := &grpcError{
		code: codes.Unimplemented,
		Msg:  fmt.Sprintf("%s plugin not registered", plugin),
	}
	return status.Error(codes.Unimplemented, ge.Error())
}
