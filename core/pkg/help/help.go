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

package help

import (
	"fmt"
)

const NITRIC_ISSUES_URL = "https://github.com/nitrictech/nitric/issues"

// standard help text to return in errors that could only be caused by a bug in Nitric.
func BugInNitricHelpText() string {
	return fmt.Sprintf("This is a bug in Nitric, please seek help: %s", NITRIC_ISSUES_URL)
}
