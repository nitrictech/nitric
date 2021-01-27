package ifaces

import (
	"context"

	"firebase.google.com/go/auth"
)

type FirebaseAuth interface {
	TenantManager() FirebaseTenantManager
}

type FirebaseTenantManager interface {
	Tenants(context.Context, string) TenantIterator
	AuthForTenant(tenantId string) (FirebaseTenantClient, error)
	CreateTenant(context.Context, *auth.TenantToCreate) (*auth.Tenant, error)
}

type TenantIterator interface {
	Next() (*auth.Tenant, error)
}

type FirebaseTenantClient interface {
	CreateUser(context.Context, *auth.UserToCreate) (*auth.UserRecord, error)
}
