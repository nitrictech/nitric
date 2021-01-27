package mocks

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"google.golang.org/api/iterator"
)

type MockFirebaseAuth struct {
	Tenants []*MockTenant
}

func (s *MockFirebaseAuth) Clear() {
	s.Tenants = make([]*MockTenant, 0)
}

func (s *MockFirebaseAuth) GetTenant(id string) *MockTenant {
	for _, t := range s.Tenants {
		if t.T.ID == id {
			return t
		}
	}

	return nil
}

func (s *MockFirebaseAuth) GetUser(tid string, uid string) *auth.UserRecord {
	for _, t := range s.Tenants {
		if t.T.ID == tid {
			for _, u := range t.Users {
				if u.UID == uid {
					return u
				}
			}
		}
	}

	return nil
}

func (s *MockFirebaseAuth) SetTenants(tenants []*MockTenant) {
	s.Tenants = tenants
}

func (s *MockFirebaseAuth) TenantManager() ifaces.FirebaseTenantManager {
	return &MockTenantManager{
		client: s,
	}
}

func NewMockFirebaseAuth() *MockFirebaseAuth {
	return &MockFirebaseAuth{
		Tenants: make([]*MockTenant, 0),
	}
}

type MockTenantManager struct {
	client *MockFirebaseAuth
}

func (s *MockTenantManager) Tenants(ctx context.Context, nextPageToken string) ifaces.TenantIterator {
	return &MockTenantIterator{
		client: s.client,
	}
}

func (s *MockTenantManager) AuthForTenant(tenantId string) (ifaces.FirebaseTenantClient, error) {
	return &MockTenantClient{
		id:     tenantId,
		client: s.client,
	}, nil
}

func (s *MockTenantManager) CreateTenant(ctx context.Context, ttc *ifaces.TenantToCreate) (*auth.Tenant, error) {
	if s.client.Tenants == nil {
		s.client.Tenants = make([]*MockTenant, 0)
	}

	for _, et := range s.client.Tenants {
		if et.T.DisplayName == ttc.DisplayName {
			return nil, fmt.Errorf("Tenant already exists")
		}
	}

	newTenant := &auth.Tenant{
		ID:                    ttc.DisplayName,
		DisplayName:           ttc.DisplayName,
		AllowPasswordSignUp:   ttc.AllowPasswordSignUp,
		EnableEmailLinkSignIn: ttc.EnableEmailLinkSignIn,
	}

	s.client.Tenants = append(s.client.Tenants, &MockTenant{
		T:     newTenant,
		Users: make([]*auth.UserRecord, 0),
	})

	return newTenant, nil
}

type MockTenantIterator struct {
	idx    int
	client *MockFirebaseAuth
}

func (m *MockTenantIterator) Next() (*auth.Tenant, error) {
	if m.idx < len(m.client.Tenants) {
		m.idx++
		mockTenant := m.client.Tenants[m.idx-1]

		return &auth.Tenant{
			ID:                    mockTenant.T.ID,
			DisplayName:           mockTenant.T.DisplayName,
			AllowPasswordSignUp:   mockTenant.T.AllowPasswordSignUp,
			EnableEmailLinkSignIn: mockTenant.T.EnableEmailLinkSignIn,
		}, nil
	}

	return nil, iterator.Done
}

type MockTenantClient struct {
	id     string
	client *MockFirebaseAuth
}

func (s *MockTenantClient) CreateUser(ctx context.Context, utc *ifaces.UserToCreate) (*auth.UserRecord, error) {
	t := s.client.GetTenant(s.id)

	if t == nil {
		return nil, fmt.Errorf("tenant does not exist")
	}

	if t.Users == nil {
		t.Users = make([]*auth.UserRecord, 0)
	}

	for _, u := range t.Users {
		if utc.UID == u.UserInfo.UID {
			return nil, fmt.Errorf("user with ID %s already exists", utc.UID)
		}

		if utc.Email == u.UserInfo.Email {
			return nil, fmt.Errorf("user with email %s already exists", utc.Email)
		}
	}

	newUser := &auth.UserRecord{
		TenantID: t.T.ID,
		Disabled: utc.Disabled,
		UserInfo: &auth.UserInfo{
			UID:         utc.UID,
			DisplayName: utc.DisplayName,
			Email:       utc.Email,
			PhoneNumber: utc.PhoneNumber,
			PhotoURL:    utc.PhotoURL,
		},
	}

	t.Users = append(t.Users, newUser)

	return newUser, nil
}

type MockTenant struct {
	T     *auth.Tenant
	Users []*auth.UserRecord
}
