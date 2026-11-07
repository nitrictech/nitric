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

func TestConfigFromAttributes_CdnSku(t *testing.T) {
	cdnAttrs := func(sku string) map[string]interface{} {
		return map[string]interface{}{
			"cdn": map[string]interface{}{"sku": sku},
		}
	}

	tests := []struct {
		name        string
		attrs       map[string]interface{}
		wantSkuName string
		wantErr     bool
	}{
		{
			name:        "defaults to standard when not specified",
			attrs:       map[string]interface{}{},
			wantSkuName: "Standard_AzureFrontDoor",
		},
		{
			name: "defaults to standard when only a custom domain is configured",
			attrs: map[string]interface{}{
				"cdn": map[string]interface{}{"domain-name": "app.example.com"},
			},
			wantSkuName: "Standard_AzureFrontDoor",
		},
		{name: "accepts standard", attrs: cdnAttrs("standard"), wantSkuName: "Standard_AzureFrontDoor"},
		{name: "accepts premium", attrs: cdnAttrs("premium"), wantSkuName: "Premium_AzureFrontDoor"},
		{name: "rejects the azure sku name", attrs: cdnAttrs("Premium_AzureFrontDoor"), wantErr: true},
		{name: "rejects unknown values", attrs: cdnAttrs("enterprise"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ConfigFromAttributes(tt.attrs)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got config with sku %q", cfg.CdnDomain.Sku)
				}

				if !strings.Contains(err.Error(), "invalid cdn sku") {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := cfg.CdnDomain.FrontDoorSkuName(); got != tt.wantSkuName {
				t.Errorf("got sku name %q, want %q", got, tt.wantSkuName)
			}
		})
	}
}
