// Copyright 2021 Nitric Technologies Pty Ltd.
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

package errors

import (
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
)

type PluginError struct {
	Code   codes.Code
	Msg    string
	Cause  error
	Plugin string
	Args   map[string]interface{}
}

func (p *PluginError) Unwrap() error {
	return p.Cause
}

func (p *PluginError) Error() string {
	if p.Cause != nil {
		// If the wrapped error is an ApiError than these should unwrap
		return fmt.Sprintf("%s: \n %s", p.Msg, p.Cause.Error())
	}

	return fmt.Sprintf("%s", p.Msg)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// Code - returns a nitric api error code from an error or Unknown if the error was not a nitric api error
func Code(e error) codes.Code {
	var pe *PluginError

	if errors.As(e, &pe) {
		return pe.Code
	}

	return codes.Unknown
}

type ErrorFactory = func(c codes.Code, msg string, cause error) error

// ErrorsWithScope - Returns a new reusable error factory with the given scope
func ErrorsWithScope(scope string, args map[string]interface{}) ErrorFactory {
	return func(code codes.Code, msg string, cause error) error {
		return &PluginError{
			Code:   code,
			Msg:    msg,
			Cause:  cause,
			Plugin: scope,
			Args:   args,
		}
	}
}
