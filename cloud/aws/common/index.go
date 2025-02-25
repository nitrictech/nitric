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

// AwsResourceName - Provides a type hint for the mapping of Nitric resource names to AWS resource names
type AwsResourceName = string

// AwsResourceArn - Provides a type hint for the mapping of Nitric resource names to AWS resource ARNs
type AwsResourceArn = string

type ApiGateway struct {
	Arn      string `json:"arn"`
	Endpoint string `json:"endpoint"`
}

type Topic struct {
	Arn             string `json:"arn"`
	StateMachineArn string `json:"stateMachineArn"`
}

// ResourceIndex - The resource index for a nitric stack
type ResourceIndex struct {
	Buckets        map[string]AwsResourceArn `json:"buckets"`
	Topics         map[string]Topic          `json:"topics"`
	KvStores       map[string]AwsResourceArn `json:"kvStores"`
	Queues         map[string]AwsResourceArn `json:"queues"`
	Secrets        map[string]AwsResourceArn `json:"secrets"`
	Apis           map[string]ApiGateway     `json:"apis"`
	HttpProxies    map[string]ApiGateway     `json:"httpProxies"`
	Websockets     map[string]ApiGateway     `json:"websockets"`
	Schedules      map[string]AwsResourceArn `json:"schedules"`
	Distributions  map[string]AwsResourceArn `json:"distributions"`
	WebsiteBuckets map[string]AwsResourceArn `json:"websiteBuckets"`
}

func NewResourceIndex() *ResourceIndex {
	return &ResourceIndex{
		Buckets:        make(map[string]AwsResourceName),
		Topics:         make(map[string]Topic),
		KvStores:       make(map[string]AwsResourceArn),
		Queues:         make(map[string]AwsResourceArn),
		Secrets:        make(map[string]AwsResourceArn),
		Apis:           make(map[string]ApiGateway),
		HttpProxies:    make(map[string]ApiGateway),
		Websockets:     make(map[string]ApiGateway),
		Schedules:      make(map[string]AwsResourceArn),
		Distributions:  make(map[string]AwsResourceArn),
		WebsiteBuckets: make(map[string]AwsResourceArn),
	}
}
