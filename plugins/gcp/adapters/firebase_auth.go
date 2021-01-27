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
	tenantClient  struct{ client *auth.TenantClient }
)

func (s firebaseAuth) TenantManager() ifaces.FirebaseTenantManager {
	return AdaptTenantManager(s.client.TenantManager)
}

func (s tenantManager) CreateTenant(ctx context.Context, t *ifaces.TenantToCreate) (*auth.Tenant, error) {
	ttc := (&auth.TenantToCreate{}).DisplayName(t.DisplayName).AllowPasswordSignUp(t.AllowPasswordSignUp).EnableEmailLinkSignIn(t.EnableEmailLinkSignIn)
	return s.client.CreateTenant(ctx, ttc)
}

func (s tenantManager) Tenants(ctx context.Context, nextPageToken string) ifaces.TenantIterator {
	return s.client.Tenants(ctx, nextPageToken)
}

func (s tenantManager) AuthForTenant(tenantId string) (ifaces.FirebaseTenantClient, error) {
	c, err := s.client.AuthForTenant(tenantId)

	if err != nil {
		return nil, err
	}

	return tenantClient{client: c}, nil
}

func (s tenantClient) CreateUser(ctx context.Context, u *ifaces.UserToCreate) (*auth.UserRecord, error) {
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
