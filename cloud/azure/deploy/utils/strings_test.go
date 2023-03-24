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
	"testing"
)

func TestJoinCamelCase(t *testing.T) {
	tests := []struct {
		name string
		ss   []string
		want string
	}{
		{
			name: "one",
			ss:   []string{"one"},
			want: "one",
		},
		{
			name: "two",
			ss:   []string{"one", "two"},
			want: "oneTwo",
		},
		{
			name: "lots",
			ss:   []string{"one", "2", "x", "eight"},
			want: "one2XEight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinCamelCase(tt.ss); got != tt.want {
				t.Errorf("joinCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringTrunc(t *testing.T) {
	tests := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{
			name: "less than",
			s:    "1234567890",
			max:  20,
			want: "1234567890",
		},
		{
			name: "max len",
			s:    "1234567890",
			max:  10,
			want: "1234567890",
		},
		{
			name: "trunc",
			s:    "1234567890",
			max:  7,
			want: "1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringTrunc(tt.s, tt.max); got != tt.want {
				t.Errorf("StringTrunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
