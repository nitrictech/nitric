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

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth interface {
	TenantManager() FirebaseTenantManager
}

type TenantToCreate struct {
	DisplayName           string
	AllowPasswordSignUp   bool
	EnableEmailLinkSignIn bool
}

type UserToCreate struct {
	Email         string
	DisplayName   string
	UID           string
	Password      string
	Disabled      bool
	EmailVerified bool
	PhoneNumber   string
	PhotoURL      string
}

type FirebaseTenantManager interface {
	Tenants(context.Context, string) TenantIterator
	AuthForTenant(tenantId string) (FirebaseTenantClient, error)
	CreateTenant(context.Context, *TenantToCreate) (*auth.Tenant, error)
}

type TenantIterator interface {
	Next() (*auth.Tenant, error)
}

type FirebaseTenantClient interface {
	CreateUser(context.Context, *UserToCreate) (*auth.UserRecord, error)
}
