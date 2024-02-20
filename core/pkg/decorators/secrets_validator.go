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

package decorators

import (
	"context"
	"fmt"

	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SecretServerValidator struct {
	inner secretspb.SecretManagerServer
}

var _ secretspb.SecretManagerServer = &SecretServerValidator{}

func validatePutRequest(req *secretspb.SecretPutRequest) error {
	if req.Secret == nil {
		return fmt.Errorf("secret cannot be nil")
	}

	if len(req.Secret.GetName()) == 0 {
		return fmt.Errorf("secret name cannot be blank")
	}

	if req.GetValue() == nil || len(req.GetValue()) == 0 {
		return fmt.Errorf("secret value cannot be blank")
	}

	return nil
}

func validateAccessRequest(req *secretspb.SecretAccessRequest) error {
	if req.SecretVersion == nil {
		return fmt.Errorf("secret version cannot be nil")
	}

	if req.SecretVersion.Secret == nil {
		return fmt.Errorf("secret cannot be nil")
	}

	if len(req.GetSecretVersion().Secret.Name) == 0 {
		return fmt.Errorf("secret name cannot be blank")
	}

	if len(req.GetSecretVersion().Version) == 0 {
		return status.Errorf(codes.InvalidArgument, "secret version cannot be blank")
	}

	return nil
}

func (s *SecretServerValidator) Put(ctx context.Context, req *secretspb.SecretPutRequest) (*secretspb.SecretPutResponse, error) {
	if err := validatePutRequest(req); err != nil {
		return nil, err
	}

	return s.inner.Put(ctx, req)
}

func (s *SecretServerValidator) Access(ctx context.Context, req *secretspb.SecretAccessRequest) (*secretspb.SecretAccessResponse, error) {
	if err := validateAccessRequest(req); err != nil {
		return nil, err
	}

	return s.inner.Access(ctx, req)
}

func SecretsServerWithValidation(inner secretspb.SecretManagerServer) *SecretServerValidator {
	return &SecretServerValidator{
		inner: inner,
	}
}
