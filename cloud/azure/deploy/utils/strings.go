// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
)

func StringTrunc(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max]
}

func JoinCamelCase(ss []string) string {
	res := ss[0]

	for i := 1; i < len(ss); i++ {
		word := ss[i]
		res += string(bytes.ToUpper([]byte{word[0]}))
		res += word[1:]
	}

	return res
}
