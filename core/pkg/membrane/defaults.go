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

package membrane

import (
	"log"
	"os"
	"strings"

	"github.com/nitrictech/nitric/core/pkg/utils"
)

func DefaultMembraneOptions() *MembraneOptions {
	options := &MembraneOptions{}

	if len(os.Args) > 1 && len(os.Args[1:]) > 0 {
		options.ChildCommand = os.Args[1:]
	} else {
		options.ChildCommand = strings.Fields(utils.GetEnv("INVOKE", ""))
		if len(options.ChildCommand) > 0 {
			log.Default().Println("Warning: use of INVOKE environment variable is deprecated and may be removed in a future version")
		}
	}

	return options
}

func fileExists(fn string) bool {
	_, err := os.Stat(fn)
	return err == nil
}
