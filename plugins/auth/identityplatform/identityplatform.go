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
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"

	"github.com/nitric-dev/membrane/sdk"
)

// IdentityPlatformAuthService - GCP Identity Platform implementation of the Nitric Auth plugin interface
type IdentityPlatformAuthService struct {
	sdk.UnimplementedAuthPlugin
	admin FirebaseAuth
}

// Get the tenant id for a given tenant
func (s *IdentityPlatformAuthService) findOrCreateTenant(tenant string) (*string, error) {
	ctx := context.Background()

	// Search for Tenant by display name first...
	tIter := s.admin.TenantManager().Tenants(ctx, "")

	for {
		if t, err := tIter.Next(); err == iterator.Done {
			// Break the loop
			// Could not find our tenant
			// So we will attempt to create one...
			break
		} else if err != nil {
			return nil, err
		} else if t.DisplayName == tenant {
			return &t.ID, nil
		}
	}

	t, err := s.admin.TenantManager().CreateTenant(context.Background(), &TenantToCreate{
		DisplayName:         tenant,
		AllowPasswordSignUp: true,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create new tenant, %v", err)
	}

	return &t.ID, nil
}

// CreateUser - Creates a new user in GCP Identity Platform (using Firebase Auth)
func (s *IdentityPlatformAuthService) Create(tenant string, id string, email string, password string) error {
	ctx := context.Background()
	tID, err := s.findOrCreateTenant(tenant)

	if err != nil {
		return err
	}

	tClient, err := s.admin.TenantManager().AuthForTenant(*tID)

	if err != nil {
		return fmt.Errorf("There was an error authorizing the requested tenant %v", err)
	}

	_, err = tClient.CreateUser(ctx, &UserToCreate{
		Email:       email,
		DisplayName: email,
		UID:         id,
		Password:    password,
	})

	if err != nil {
		return fmt.Errorf("There was an error creating the new user: %v", err)
	}

	return nil
}

// New - Creates a new instance of the Identity Platform auth plugin
func New() (sdk.UserService, error) {
	ctx := context.Background()

	//credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	//if credentialsError != nil {
	//	return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	//}

	//credOpt := option.WithCredentialsJSON(credentials.JSON)
	app, err := firebase.NewApp(ctx, nil)

	//fmt.Println("Creds:", credOpt)

	if err != nil {
		return nil, fmt.Errorf("Error instansiating Firebase App credentials %v", err)
	}

	// fmt.Println("App:", app)

	authClient, err := app.Auth(ctx)

	if err != nil {
		return nil, fmt.Errorf("Error instansiating firebase auth client %v", err)
	}

	return &IdentityPlatformAuthService{
		admin: AdaptFirebaseAuth(authClient),
	}, nil
}

func NewWithClient(client FirebaseAuth) sdk.UserService {
	return &IdentityPlatformAuthService{
		admin: client,
	}
}
