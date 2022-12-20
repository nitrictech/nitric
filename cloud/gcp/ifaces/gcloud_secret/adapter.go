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

package ifaces_gcloud_secret

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	gax "github.com/googleapis/gax-go/v2"
)

type realClient struct {
	*secretmanager.Client
}

func NewClient(ctx context.Context) (SecretManagerClient, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &realClient{Client: c}, nil
}

func (r *realClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, co ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return r.Client.AccessSecretVersion(ctx, req, co...)
}

func (r *realClient) AddSecretVersion(ctx context.Context, req *secretmanagerpb.AddSecretVersionRequest, co ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return r.Client.AddSecretVersion(ctx, req, co...)
}

func (r *realClient) UpdateSecret(ctx context.Context, req *secretmanagerpb.UpdateSecretRequest, co ...gax.CallOption) (*secretmanagerpb.Secret, error) {
	return r.Client.UpdateSecret(ctx, req, co...)
}

func (r *realClient) ListSecrets(ctx context.Context, req *secretmanagerpb.ListSecretsRequest, co ...gax.CallOption) SecretIterator {
	return r.Client.ListSecrets(ctx, req, co...)
}
