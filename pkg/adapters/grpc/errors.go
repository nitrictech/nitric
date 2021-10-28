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
	"reflect"

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
		ed.Scope = &v1.ErrorScope{
			Service: operation,
			Plugin:  pe.Plugin,
		}
		if len(pe.Args) > 0 {
			args := make(map[string]string)
			for k, v := range pe.Args {
				args[k] = LogArg(v)
			}
			ed.Scope.Args = args
		}

		s := status.New(code, pe.Msg)
		s, _ = s.WithDetails(ed)

		return s.Err()

	} else {
		return newGrpcErrorWithCode(codes.Internal, operation, err)
	}
}

func newGrpcErrorWithCode(code codes.Code, operation string, err error) error {
	ed := &v1.ErrorDetails{}
	ed.Message = err.Error()
	ed.Scope = &v1.ErrorScope{
		Service: operation,
	}

	s := status.New(code, err.Error())
	s, _ = s.WithDetails(ed)

	return s.Err()
}

// Provides generic error for unregistered plugins
func NewPluginNotRegisteredError(plugin string) error {
	ed := &v1.ErrorDetails{}
	ed.Message = fmt.Sprintf("%s plugin not registered", plugin)

	s := status.New(codes.Unimplemented, ed.Message)
	s, _ = s.WithDetails(ed)

	return s.Err()
}

func LogArg(arg interface{}) string {
	value := getValue(arg)

	if value.Kind() == reflect.Struct {

		str := "{"
		for i := 0; i < value.NumField(); i++ {

			fieldType := value.Type().Field(i)
			tag := fieldType.Tag.Get("log")
			if tag == "" || tag == "-" {
				continue
			}

			if len(str) > 1 {
				str += ", "
			}

			field := value.Field(i)
			str += fieldType.Name + ": " + LogArg(field.Interface())
		}
		str += "}"

		return str

	} else if value.Kind() == reflect.Map {
		str := "{"

		for k, v := range arg.(map[string]interface{}) {
			if len(str) > 1 {
				str += ", "
			}
			str += fmt.Sprintf("%v", k) + ": " + LogArg(v)
		}

		str += "}"

		return str

	} else {
		return fmt.Sprintf("%v", arg)
	}
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}
