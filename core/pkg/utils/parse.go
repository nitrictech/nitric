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

import (
	"strconv"
)

// PercentFromIntString returns a float between 0.0 to 1 representing a percentage.
// this is converted from a string int in the range "0" to "100".
func PercentFromIntString(in string) (float64, error) {
	intVar, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}

	if intVar >= 100 {
		return 1, nil
	} else if intVar <= 0 {
		return 0, nil
	}

	return float64(intVar) / float64(100), nil
}
