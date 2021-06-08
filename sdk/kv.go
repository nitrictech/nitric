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

type QueryExpression struct {
	Operand  string
	Operator string
	Value    string
}

type QueryResult struct {
	Data        []map[string]interface{}
	PagingToken map[string]interface{}
}

func (e QueryExpression) String() string {
	return fmt.Sprintf("{Operand: '%v', Operator: '%v', Value: '%v'}", e.Operand, e.Operator, e.Value)
}

func (q *QueryResult) String() string {
	return fmt.Sprintf("Data len:%v, PagingToken: %v", len(q.Data), q.PagingToken)
}

// The base KeyValue Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type KeyValueService interface {
	Put(string, map[string]interface{}, map[string]interface{}) error
	Get(string, map[string]interface{}) (map[string]interface{}, error)
	Delete(string, map[string]interface{}) error
	Query(string, []QueryExpression, int, map[string]interface{}) (*QueryResult, error)
}

type UnimplementedKeyValuePlugin struct {
	KeyValueService
}

func (p *UnimplementedKeyValuePlugin) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Delete(collection string, key map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedKeyValuePlugin) Query(collection string, expressions []QueryExpression, limit int, pagingToken map[string]interface{}) (*QueryResult, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
