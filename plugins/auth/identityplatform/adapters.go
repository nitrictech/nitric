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

package identity_platform_service

import (
	"context"

	auth "firebase.google.com/go/v4/auth"
)

func AdaptFirebaseAuth(client *auth.Client) FirebaseAuth {
	return &firebaseAuth{client: client}
}

func AdaptTenantManager(client *auth.TenantManager) FirebaseTenantManager {
	return &tenantManager{client: client}
}

type (
	firebaseAuth  struct{ client *auth.Client }
	tenantManager struct{ client *auth.TenantManager }
	tenantClient  struct{ client *auth.TenantClient }
)

func (s firebaseAuth) TenantManager() FirebaseTenantManager {
	return AdaptTenantManager(s.client.TenantManager)
}

func (s tenantManager) CreateTenant(ctx context.Context, t *TenantToCreate) (*auth.Tenant, error) {
	ttc := (&auth.TenantToCreate{}).DisplayName(t.DisplayName).AllowPasswordSignUp(t.AllowPasswordSignUp).EnableEmailLinkSignIn(t.EnableEmailLinkSignIn)
	return s.client.CreateTenant(ctx, ttc)
}

func (s tenantManager) Tenants(ctx context.Context, nextPageToken string) TenantIterator {
	return s.client.Tenants(ctx, nextPageToken)
}

func (s tenantManager) AuthForTenant(tenantId string) (FirebaseTenantClient, error) {
	c, err := s.client.AuthForTenant(tenantId)

	if err != nil {
		return nil, err
	}

	return tenantClient{client: c}, nil
}

func (s tenantClient) CreateUser(ctx context.Context, u *UserToCreate) (*auth.UserRecord, error) {
	utc := &auth.UserToCreate{}

	utc.DisplayName(u.DisplayName)
	utc.Disabled(u.Disabled)
	utc.Email(u.Email)
	utc.EmailVerified(u.EmailVerified)
	utc.Password(u.Password)
	utc.PhoneNumber(u.PhoneNumber)
	utc.PhotoURL(u.PhotoURL)
	utc.UID(u.UID)

	return s.client.CreateUser(ctx, utc)
}
