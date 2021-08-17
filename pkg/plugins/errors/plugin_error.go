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
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
)

type PluginError struct {
	code  codes.Code
	msg   string
	cause error
}

func (p *PluginError) Unwrap() error {
	return p.cause
}

func (p *PluginError) Error() string {
	if p.cause != nil {
		// If the wrapped error is an ApiError than these should unwrap
		return fmt.Sprintf("%s: \n %s", p.msg, p.cause.Error())
	}

	return fmt.Sprintf("%s", p.msg)
}

// Code - returns a nitric api error code from an error or Unknown if the error was not a nitric api error
func Code(e error) codes.Code {
	if pe, ok := e.(*PluginError); ok {
		return pe.code
	}

	return codes.Unknown
}

// New - Creates a new nitric API error
func new(c codes.Code, msg string) error {
	return &PluginError{
		code: c,
		msg:  msg,
	}
}

// NewWithCause - Creates a new nitric API error with the given error as it's cause
func newWithCause(c codes.Code, msg string, cause error) error {
	return &PluginError{
		code:  c,
		msg:   msg,
		cause: cause,
	}
}

// ErrorsWithScope - Returns a new reusable error factory with the given scope
func ErrorsWithScope(s string, ctx ...interface{}) func(c codes.Code, msg string, cause error) error {
	return func(c codes.Code, msg string, cause error) error {
		sMsg := fmt.Sprintf("%s(%v): %s", s, ctx, msg)

		if cause == nil {
			return new(c, sMsg)
		}

		return newWithCause(c, sMsg, cause)
	}
}
