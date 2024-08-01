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

package provider

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type ErrorHandler = func(err error) error

func WithErrorHandler(handler ErrorHandler) func(*PulumiProviderServer) {
	return func(s *PulumiProviderServer) {
		s.errorHandlers = append(s.errorHandlers, handler)
	}
}

func handleCommonErrors(err error) error {
	// Check for common Pulumi 'autoError' types
	if auto.IsConcurrentUpdateError(err) {
		if pe := parsePulumiError(err); pe != nil {
			err = pe
		}
		err = fmt.Errorf("the pulumi stack file is locked.\nThis occurs when a previous deployment is still in progress or was interrupted.\n%w", err)
	} else if auto.IsSelectStack404Error(err) {
		err = fmt.Errorf("stack not found. %w", err)
	} else if auto.IsCreateStack409Error(err) {
		err = fmt.Errorf("failed to create Pulumi stack, this may be a bug in nitric. Seek help https://github.com/nitrictech/nitric/issues\n%w", err)
	} else if auto.IsCompilationError(err) {
		err = fmt.Errorf("failed to compile Pulumi program, this may be a bug in your chosen provider or with nitric. Seek help https://github.com/nitrictech/nitric/issues\n%w", err)
	} else if pe := parsePulumiError(err); pe != nil {
		err = pe
	}

	return err
}
