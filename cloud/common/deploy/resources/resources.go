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

package resources

type ResourceType string

const (
	API         ResourceType = "api"
	Bucket      ResourceType = "bucket"
	Collection  ResourceType = "collection"
	Service     ResourceType = "service"
	HttpProxy   ResourceType = "http-proxy"
	Policy      ResourceType = "policy"
	Queue       ResourceType = "queue"
	Schedule    ResourceType = "schedule"
	Secret      ResourceType = "secret"
	Stack       ResourceType = "stack"
	Topic       ResourceType = "topic"
	Websocket   ResourceType = "websocket"
	SqlDatabase ResourceType = "sql-database"
)
