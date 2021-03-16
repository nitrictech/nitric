package identity_platform_service

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/nitric-dev/membrane/plugins/gcp/adapters"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

// IdentityPlatformAuthService - GCP Identity Platform implementation of the Nitric Auth plugin interface
type IdentityPlatformAuthService struct {
	sdk.UnimplementedAuthPlugin
	admin ifaces.FirebaseAuth
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

	t, err := s.admin.TenantManager().CreateTenant(context.Background(), &ifaces.TenantToCreate{
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

	_, err = tClient.CreateUser(ctx, &ifaces.UserToCreate{
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

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	credOpt := option.WithCredentialsJSON(credentials.JSON)

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: credentials.ProjectID,
	}, credOpt)

	if err != nil {
		return nil, fmt.Errorf("Error instansiating Firebase App credentials")
	}

	authClient, err := app.Auth(ctx)

	if err != nil {
		return nil, fmt.Errorf("Error instansiating firebase auth client")
	}

	return &IdentityPlatformAuthService{
		admin: adapters.AdaptFirebaseAuth(authClient),
	}, nil
}

func NewWithClient(client ifaces.FirebaseAuth) sdk.UserService {
	return &IdentityPlatformAuthService{
		admin: client,
	}
}
