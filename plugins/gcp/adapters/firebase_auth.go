package adapters

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
)

func AdaptFirebaseAuth(client *auth.Client) ifaces.FirebaseAuth {
	return &firebaseAuth{client: client}
}

func AdaptTenantManager(client *auth.TenantManager) ifaces.FirebaseTenantManager {
	return &tenantManager{client: client}
}

type (
	firebaseAuth  struct{ client *auth.Client }
	tenantManager struct{ client *auth.TenantManager }
)

func (s firebaseAuth) TenantManager() ifaces.FirebaseTenantManager {
	return AdaptTenantManager(s.client.TenantManager)
}

func (s tenantManager) CreateTenant(ctx context.Context, t *auth.TenantToCreate) (*auth.Tenant, error) {
	return s.client.CreateTenant(ctx, t)
}

func (s tenantManager) Tenants(ctx context.Context, nextPageToken string) ifaces.TenantIterator {
	return s.client.Tenants(ctx, nextPageToken)
}

func (s tenantManager) AuthForTenant(tenantId string) (ifaces.FirebaseTenantClient, error) {
	return s.client.AuthForTenant(tenantId)
}
