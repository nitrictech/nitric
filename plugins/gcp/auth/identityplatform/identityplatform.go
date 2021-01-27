package identity_platform_plugin

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/nitric-dev/membrane/plugins/sdk"
)

// IdentityPlatformPlugin - GCP Identity Platform implementation of the Nitric Auth plugin interface
type IdentityPlatformPlugin struct {
	sdk.UnimplementedAuthPlugin
	admin *auth.Client
}

// Get the tenant id for a given tenant
func (s *IdentityPlatformPlugin) findOrCreateTenant(tenant string) (*string, error) {
	ctx := context.Background()

	// Search for Tenant by display name first...
	tIter := s.admin.TenantManager.Tenants(ctx, "")

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
		} else {
			fmt.Println("Comparing %s -> %s", tenant, t.DisplayName)
		}
	}

	tCreate := (&auth.TenantToCreate{}).DisplayName(tenant).AllowPasswordSignUp(true)

	t, err := s.admin.TenantManager.CreateTenant(context.Background(), tCreate)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new tenant, %v", err)
	}

	return &t.ID, nil
}

// CreateUser - Creates a new user in GCP Identity Platform (using Firebase Auth)
func (s *IdentityPlatformPlugin) CreateUser(tenant string, id string, email string, password string) error {
	ctx := context.Background()

	tID, err := s.findOrCreateTenant(tenant)

	if err != nil {
		return err
	}

	tClient, err := s.admin.TenantManager.AuthForTenant(*tID)

	if err != nil {
		return fmt.Errorf("There was an error authorizing the requested tenant %v", err)
	}

	uCreate := (&auth.UserToCreate{}).Email(email).DisplayName(email).UID(id).Password(password)

	_, err = tClient.CreateUser(ctx, uCreate)

	if err != nil {
		return fmt.Errorf("There was an error creating the new user: %v", err)
	}

	return nil
}

// New - Creates a new instance of the Cognito auth plugin
func New() (sdk.AuthPlugin, error) {
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

	return &IdentityPlatformPlugin{
		admin: authClient,
	}, nil
}
