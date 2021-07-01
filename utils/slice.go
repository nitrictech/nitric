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

package utils

// Return index of the test in the given slice values, or -1 if not found
func IndexOf(values []string, test string) int {
	if len(values) == 0 {
		return -1
	}
	for index, val := range values {
		if val == test {
			return index
		}
	}
	return -1
}

// Return a string slice with the specified index removed (zero indexed).
// The returned slice will maintain the original order.
// This function  will return the original slice if the remove operation is not valid
func Remove(slice []string, index int) []string {
	if len(slice) == 0 || index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
