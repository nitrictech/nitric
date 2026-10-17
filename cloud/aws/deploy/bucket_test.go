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

package deploy

import (
	"sync"
	"testing"

	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type bucketMocks struct {
	mu        sync.Mutex
	resources []mockResource
}

type mockResource struct {
	TypeToken string
	Name      string
	Inputs    resource.PropertyMap
}

func (m *bucketMocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.resources = append(m.resources, mockResource{
		TypeToken: args.TypeToken,
		Name:      args.Name,
		Inputs:    args.Inputs,
	})

	// Return a fake ID and the inputs as outputs
	return args.Name + "-id", args.Inputs, nil
}

func (m *bucketMocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func (m *bucketMocks) findResources(typeToken string) []mockResource {
	m.mu.Lock()
	defer m.mu.Unlock()

	var found []mockResource
	for _, r := range m.resources {
		if r.TypeToken == typeToken {
			found = append(found, r)
		}
	}
	return found
}

func TestBucket_WithCorsRules_CreatesS3CorsConfiguration(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider := &NitricAwsPulumiProvider{
			StackId:             "test-stack",
			AwsConfig:           &common.AwsConfig{},
			Buckets:             map[string]*s3.Bucket{},
			BucketNotifications: map[string]*s3.BucketNotification{},
		}

		return provider.Bucket(ctx, nil, "test-bucket", &deploymentspb.Bucket{
			CorsRules: []*deploymentspb.BucketCorsRule{
				{
					AllowedOrigins: []string{"https://example.com"},
					AllowedMethods: []string{"GET", "POST"},
					AllowedHeaders: []string{"Content-Type", "Authorization"},
					ExposeHeaders:  []string{"ETag"},
					MaxAgeSeconds:  3600,
				},
			},
		})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify an S3 bucket was created
	buckets := mocks.findResources("aws:s3/bucket:Bucket")
	if len(buckets) != 1 {
		t.Fatalf("expected 1 S3 bucket, got %d", len(buckets))
	}

	// Verify a CORS configuration was created
	corsConfigs := mocks.findResources("aws:s3/bucketCorsConfigurationV2:BucketCorsConfigurationV2")
	if len(corsConfigs) != 1 {
		t.Fatalf("expected 1 CORS configuration, got %d", len(corsConfigs))
	}

	// Verify the CORS rule field values
	corsRules, ok := corsConfigs[0].Inputs["corsRules"]
	if !ok {
		t.Fatal("expected CORS configuration to have 'corsRules' input")
	}
	rules := corsRules.ArrayValue()
	if len(rules) != 1 {
		t.Fatalf("expected 1 CORS rule, got %d", len(rules))
	}
	rule := rules[0].ObjectValue()

	origins := rule["allowedOrigins"].ArrayValue()
	if len(origins) != 1 || origins[0].StringValue() != "https://example.com" {
		t.Fatalf("unexpected allowedOrigins: %v", origins)
	}
	methods := rule["allowedMethods"].ArrayValue()
	if len(methods) != 2 || methods[0].StringValue() != "GET" {
		t.Fatalf("unexpected allowedMethods: %v", methods)
	}
	headers := rule["allowedHeaders"].ArrayValue()
	if len(headers) != 2 || headers[0].StringValue() != "Content-Type" {
		t.Fatalf("unexpected allowedHeaders: %v", headers)
	}
	expose := rule["exposeHeaders"].ArrayValue()
	if len(expose) != 1 || expose[0].StringValue() != "ETag" {
		t.Fatalf("unexpected exposeHeaders: %v", expose)
	}
	maxAge, ok := rule["maxAgeSeconds"]
	if !ok || maxAge.NumberValue() != 3600 {
		t.Fatalf("unexpected maxAgeSeconds: %v", maxAge)
	}
}

func TestBucket_WithoutCorsRules_NoCorsConfiguration(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider := &NitricAwsPulumiProvider{
			StackId:             "test-stack",
			AwsConfig:           &common.AwsConfig{},
			Buckets:             map[string]*s3.Bucket{},
			BucketNotifications: map[string]*s3.BucketNotification{},
		}

		return provider.Bucket(ctx, nil, "test-bucket", &deploymentspb.Bucket{})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify no CORS configuration was created
	corsConfigs := mocks.findResources("aws:s3/bucketCorsConfigurationV2:BucketCorsConfigurationV2")
	if len(corsConfigs) != 0 {
		t.Fatalf("expected 0 CORS configurations, got %d", len(corsConfigs))
	}
}

func TestBucket_MultipleCorsRules(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider := &NitricAwsPulumiProvider{
			StackId:             "test-stack",
			AwsConfig:           &common.AwsConfig{},
			Buckets:             map[string]*s3.Bucket{},
			BucketNotifications: map[string]*s3.BucketNotification{},
		}

		return provider.Bucket(ctx, nil, "test-bucket", &deploymentspb.Bucket{
			CorsRules: []*deploymentspb.BucketCorsRule{
				{
					AllowedOrigins: []string{"https://example.com"},
					AllowedMethods: []string{"GET"},
					MaxAgeSeconds:  300,
				},
				{
					AllowedOrigins: []string{"https://other.com"},
					AllowedMethods: []string{"PUT", "POST"},
					AllowedHeaders: []string{"*"},
					MaxAgeSeconds:  600,
				},
			},
		})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify exactly one CORS configuration resource (containing both rules)
	corsConfigs := mocks.findResources("aws:s3/bucketCorsConfigurationV2:BucketCorsConfigurationV2")
	if len(corsConfigs) != 1 {
		t.Fatalf("expected 1 CORS configuration, got %d", len(corsConfigs))
	}
}
