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

	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
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

func newTestGcpProvider(ctx *pulumi.Context) (*NitricGcpPulumiProvider, error) {
	project := &Project{Name: "test-project"}
	err := ctx.RegisterComponentResource("ntiricgcp:project:GcpProject", "test-project", project)
	if err != nil {
		return nil, err
	}

	return &NitricGcpPulumiProvider{
		CommonStackDetails: &deploy.CommonStackDetails{Region: "us-central1"},
		Buckets:            map[string]*storage.Bucket{},
		StackId:            "test-stack",
		Project:            project,
	}, nil
}

func TestBucket_WithCorsRules_SetsCorsOnBucket(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider, err := newTestGcpProvider(ctx)
		if err != nil {
			return err
		}

		return provider.Bucket(ctx, nil, "test-bucket", &deploymentspb.Bucket{
			CorsRules: []*deploymentspb.BucketCorsRule{
				{
					AllowedOrigins: []string{"https://example.com"},
					AllowedMethods: []string{"GET", "POST"},
					AllowedHeaders: []string{"Content-Type"},
					ExposeHeaders:  []string{"ETag"},
					MaxAgeSeconds:  3600,
				},
			},
		})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// GCP sets CORS on the bucket itself (not a separate resource)
	buckets := mocks.findResources("gcp:storage/bucket:Bucket")
	if len(buckets) != 1 {
		t.Fatalf("expected 1 GCS bucket, got %d", len(buckets))
	}

	// Verify the bucket has cors property set
	bucket := buckets[0]
	cors, ok := bucket.Inputs["cors"]
	if !ok {
		t.Fatal("expected bucket to have 'cors' property set")
	}

	corsArr := cors.ArrayValue()
	if len(corsArr) != 1 {
		t.Fatalf("expected 1 CORS rule, got %d", len(corsArr))
	}

	// Verify field values
	rule := corsArr[0].ObjectValue()

	origins := rule["origins"].ArrayValue()
	if len(origins) != 1 || origins[0].StringValue() != "https://example.com" {
		t.Fatalf("unexpected origins: %v", origins)
	}
	methods := rule["methods"].ArrayValue()
	if len(methods) != 2 || methods[0].StringValue() != "GET" {
		t.Fatalf("unexpected methods: %v", methods)
	}
	// GCP maps ExposeHeaders to ResponseHeaders
	responseHeaders := rule["responseHeaders"].ArrayValue()
	if len(responseHeaders) != 1 || responseHeaders[0].StringValue() != "ETag" {
		t.Fatalf("unexpected responseHeaders: %v", responseHeaders)
	}
	maxAge, ok := rule["maxAgeSeconds"]
	if !ok || maxAge.NumberValue() != 3600 {
		t.Fatalf("unexpected maxAgeSeconds: %v", maxAge)
	}
}

func TestBucket_MultipleCorsRules(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider, err := newTestGcpProvider(ctx)
		if err != nil {
			return err
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
					ExposeHeaders:  []string{"X-Custom"},
					MaxAgeSeconds:  600,
				},
			},
		})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buckets := mocks.findResources("gcp:storage/bucket:Bucket")
	if len(buckets) != 1 {
		t.Fatalf("expected 1 GCS bucket, got %d", len(buckets))
	}

	cors, ok := buckets[0].Inputs["cors"]
	if !ok {
		t.Fatal("expected bucket to have 'cors' property set")
	}
	if len(cors.ArrayValue()) != 2 {
		t.Fatalf("expected 2 CORS rules, got %d", len(cors.ArrayValue()))
	}
}

func TestBucket_WithoutCorsRules_NoCorsOnBucket(t *testing.T) {
	mocks := &bucketMocks{}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		provider, err := newTestGcpProvider(ctx)
		if err != nil {
			return err
		}

		return provider.Bucket(ctx, nil, "test-bucket", &deploymentspb.Bucket{})
	}, pulumi.WithMocks("test", "test", mocks))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buckets := mocks.findResources("gcp:storage/bucket:Bucket")
	if len(buckets) != 1 {
		t.Fatalf("expected 1 GCS bucket, got %d", len(buckets))
	}

	// Verify no cors property
	bucket := buckets[0]
	cors, ok := bucket.Inputs["cors"]
	if ok && len(cors.ArrayValue()) > 0 {
		t.Fatal("expected bucket to NOT have 'cors' property set")
	}
}
