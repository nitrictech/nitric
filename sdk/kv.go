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

package sdk

import "fmt"

// The base KeyValue Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type KeyValueService interface {
	Put(string, string, map[string]interface{}) error
	Get(string, string) (map[string]interface{}, error)
	Delete(string, string) error
}

type UnimplementedKeyValuePlugin struct {
	KeyValueService
}

func (p *UnimplementedKeyValuePlugin) Put(collection string, key string, value map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Get(collection string, key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Delete(collection string, key string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
