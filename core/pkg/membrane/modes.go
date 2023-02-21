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
	"fmt"
	"strings"
)

// SourceType enum
type Mode int

const (
	// Mode_Faas Facilitates FaaS via gRPC FaaS Server
	Mode_Faas Mode = iota
	// Mode_HttpProxy is designed for integration of monoliths into a nitric application
	Mode_HttpProxy
)

var modes = [...]string{"FAAS", "HTTP_PROXY"}

func (m Mode) String() string {
	return modes[m]
}

func ModeFromString(modeString string) (Mode, error) {
	for i, mode := range modes {
		if mode == modeString {
			return Mode(i), nil
		}
	}
	return -1, fmt.Errorf("Invalid mode %s, supported modes are: %s", modeString, strings.Join(modes[:], ", "))
}
