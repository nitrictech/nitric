package identity_platform_service

import (
	"context"

	"firebase.google.com/go/auth"
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
