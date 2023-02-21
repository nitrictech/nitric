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

package document

import (
	"context"
	"fmt"
)

// MaxSubCollectionDepth - maximum number of parents a collection can support.
// Depth is a count of the number of parents for a collection.
// e.g. a collection with no parent has a depth of 0
// a collection with a parent has a depth of 1
const MaxSubCollectionDepth int = 1

type Collection struct {
	Name   string `log:"Name"`
	Parent *Key   `log:"Parent"`
}

type Key struct {
	Collection *Collection `log:"Collection"`
	Id         string      `log:"Id"`
}

type Document struct {
	Key     *Key
	Content map[string]interface{}
}

type QueryExpression struct {
	Operand  string
	Operator string
	Value    interface{}
}

type QueryResult struct {
	Documents   []Document
	PagingToken map[string]string
}

type DocumentIterator = func() (*Document, error)

// The base Document Plugin interface
// Use this over proto definitions to remove dependency on protobuf in the plugin internally
// and open options to adding additional non-grpc interfaces
type DocumentService interface {
	Get(context.Context, *Key) (*Document, error)
	Set(context.Context, *Key, map[string]interface{}) error
	Delete(context.Context, *Key) error
	Query(context.Context, *Collection, []QueryExpression, int, map[string]string) (*QueryResult, error)
	QueryStream(context.Context, *Collection, []QueryExpression, int) DocumentIterator
}

type UnimplementedDocumentPlugin struct {
	DocumentService
}

func (p *UnimplementedDocumentPlugin) Get(ctx context.Context, key *Key) (*Document, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentPlugin) Set(ctx context.Context, key *Key, content map[string]interface{}) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentPlugin) Delete(ctx context.Context, key *Key) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentPlugin) Query(ctx context.Context, collection *Collection, expressions []QueryExpression, limit int, pagingToken map[string]string) (*QueryResult, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (p *UnimplementedDocumentPlugin) QueryStream(ctx context.Context, collection *Collection, expressions []QueryExpression, limit int) DocumentIterator {
	return func() (*Document, error) {
		return nil, fmt.Errorf("UNIMPLEMENTED")
	}
}
