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
