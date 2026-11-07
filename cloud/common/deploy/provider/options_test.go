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
	"errors"
	"strings"
	"testing"
)

func TestExplainCommonErrs(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantContains []string
		wantExact    string
	}{
		{
			name: "explains a missing state backend in a non-interactive session",
			err: errors.New("failed to select stack: exit status 3\ncode: 3\nstdout: \n" +
				"stderr: error: PULUMI_ACCESS_TOKEN must be set for login during non-interactive CLI sessions\n"),
			wantContains: []string{
				"no Pulumi state backend is configured",
				"PULUMI_BACKEND_URL",
				"https://nitric.io/docs/providers/pulumi#state-backends",
				// the original error is kept for troubleshooting
				"failed to select stack: exit status 3",
				"PULUMI_ACCESS_TOKEN must be set for login during non-interactive CLI sessions",
			},
		},
		{
			name:      "does not explain unrelated errors",
			err:       errors.New("something went wrong"),
			wantExact: "something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := explainCommonErrs(tt.err).Error()

			if tt.wantExact != "" && got != tt.wantExact {
				t.Errorf("got %q, want %q", got, tt.wantExact)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("got %q, want it to contain %q", got, want)
				}
			}
		})
	}
}
