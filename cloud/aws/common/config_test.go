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

package common

import (
	"strings"
	"testing"
)

func TestConfigFromAttributes_ResourceResolver(t *testing.T) {
	tests := []struct {
		name        string
		attrs       map[string]interface{}
		wantValue   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "defaults to ssm when not specified",
			attrs:     map[string]interface{}{},
			wantValue: ResourceResolverSSM,
		},
		{
			name:      "accepts ssm explicitly",
			attrs:     map[string]interface{}{"resource-resolver": "ssm"},
			wantValue: ResourceResolverSSM,
		},
		{
			name:      "accepts tagging",
			attrs:     map[string]interface{}{"resource-resolver": "tagging"},
			wantValue: ResourceResolverTagging,
		},
		{
			name:        "rejects invalid value",
			attrs:       map[string]interface{}{"resource-resolver": "taging"},
			wantErr:     true,
			errContains: "invalid resource-resolver",
		},
		{
			name:        "rejects arbitrary string",
			attrs:       map[string]interface{}{"resource-resolver": "something-else"},
			wantErr:     true,
			errContains: "invalid resource-resolver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ConfigFromAttributes(tt.attrs)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errContains)
				}
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Fatalf("expected error containing %q, got %q", tt.errContains, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.ResourceResolver != tt.wantValue {
				t.Errorf("ResourceResolver = %q, want %q", cfg.ResourceResolver, tt.wantValue)
			}
		})
	}
}
