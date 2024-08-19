// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gateway

import "fmt"

// An error indicating the JSON event type is not handled by this lambda gateway
type UnhandledLambdaEventError struct {
	Message string
	Cause   error
}

func (e *UnhandledLambdaEventError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

func (e *UnhandledLambdaEventError) Unwrap() error {
	return e.Cause
}

func NewUnhandledLambdaEventError(cause error) *UnhandledLambdaEventError {
	return &UnhandledLambdaEventError{
		Message: "the nitric lambda gateway does not handle this lambda event type",
		Cause:   cause,
	}
}
