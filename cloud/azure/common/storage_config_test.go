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

func TestConfigFromAttributes_StorageReplicationType(t *testing.T) {
	storageAttrs := func(replicationType string) map[string]interface{} {
		return map[string]interface{}{
			"storage": map[string]interface{}{"replication-type": replicationType},
		}
	}

	tests := []struct {
		name      string
		attrs     map[string]interface{}
		wantValue string
		wantErr   bool
	}{
		{
			name:      "defaults to LRS when not specified",
			attrs:     map[string]interface{}{},
			wantValue: "LRS",
		},
		{name: "accepts LRS", attrs: storageAttrs("LRS"), wantValue: "LRS"},
		{name: "accepts ZRS", attrs: storageAttrs("ZRS"), wantValue: "ZRS"},
		{name: "accepts GRS", attrs: storageAttrs("GRS"), wantValue: "GRS"},
		{name: "accepts RAGRS", attrs: storageAttrs("RAGRS"), wantValue: "RAGRS"},
		{name: "accepts GZRS", attrs: storageAttrs("GZRS"), wantValue: "GZRS"},
		{name: "accepts RAGZRS", attrs: storageAttrs("RAGZRS"), wantValue: "RAGZRS"},
		{name: "rejects lowercase values", attrs: storageAttrs("grs"), wantErr: true},
		{name: "rejects the azure sku name", attrs: storageAttrs("Standard_GRS"), wantErr: true},
		{name: "rejects unknown values", attrs: storageAttrs("XRS"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ConfigFromAttributes(tt.attrs)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got config with replication type %q", cfg.Storage.ReplicationType)
				}

				if !strings.Contains(err.Error(), "invalid storage replication-type") {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Storage.ReplicationType != tt.wantValue {
				t.Errorf("got replication type %q, want %q", cfg.Storage.ReplicationType, tt.wantValue)
			}
		})
	}
}
