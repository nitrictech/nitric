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

package grpc_errors

import (
	"fmt"

	// structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type ScopedErrorFactory = func(c codes.Code, msg string, cause error) error

func ErrorsWithScope(scope string) ScopedErrorFactory {
	return func(code codes.Code, msg string, cause error) error {
		scopedMsg := fmt.Sprintf("%s %s", scope, msg)

		st := status.New(code, scopedMsg)
		if cause != nil {
			errorDetails := map[string]interface{}{
				"cause": cause.Error(),
			}

			detail, err := structpb.NewStruct(errorDetails)
			if err != nil {
				return status.Error(code, fmt.Sprintf("%s - %s", scopedMsg, cause.Error()))
			}

			st, err = st.WithDetails(detail)
			if err != nil {
				return status.Error(code, fmt.Sprintf("%s - %s", scopedMsg, cause.Error()))
			}
		}

		return st.Err()
	}
}
