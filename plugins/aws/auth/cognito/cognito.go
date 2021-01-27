package cognito_plugin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// CognitoPlugin - Cognito implementation of the Nitric Auth plugin interface
type CognitoPlugin struct {
	sdk.UnimplementedAuthPlugin
	client cognitoidentityprovider.CognitoIdentityProvider
}

// Get the client id for a given user pool
func (s *CognitoPlugin) findOrCreateUserPoolForTenant(tenant string) (*string, error) {
	out, err := s.client.ListUserPoolClients(&cognitoidentityprovider.ListUserPoolsInput{})

	for _, client := range out.UserPoolClients {
		if client.ClientName == tenant {
			return client.ClientId, nil
		}
	}

	// TODO: Determine if the user pool already exists....
	// Otherwise create a new one...
	userPool, err := s.client.CreateUserPool(&cognitoidentityprovider.CreateUserPoolInput{
		PoolName: tenant,
	})

	if err != nil {
		return nil, fmt.Errorf("Error creating userpool for tenant: %s", tenant)
	}

	client, err := s.client.CreateUserPoolClient(&cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId: userPool.UserPool.Id,
		ClientName: tenant,
		ExplicitAuthFlows: []string{"ALLOW_USER_PASSWORD_AUTH"},
		GenerateSecret: false,
	})

	if err != nil {
		return nil, err
	}

	return client.UserPoolClient.ClientId, nil
}

// CreateUser - Creates a new user in AWS cognito
func (s *CognitoPlugin) CreateUser(tenant string, id string, email string, password string) error {
	// Attempt to sign up the user...
	upClient, err := s.findOrCreateUserPoolForTenant()

	if err != nil {
		return fmt.Errorf("Could not SignUp user: %v", err)
	}

	_, err := s.client.SignUp(&cognitoidentityprovider.SignUpInput{
		// TODO: Need to determine the client id in this case
		// For email/password authentication, will likely just do a single user pool for the stack...
		ClientId: "",
		Password: password,
		Username: id,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			&cognitoidentityprovider.AttributeType{
				Name: "email",
				Value: email,
			}
		},
	})

	if err != nil {
		return fmt.Errorf("There was an error signing up the user: %v", err)
	}

	return nil
}

// New - Creates a new instance of the Cognito auth plugin
func New() (sdk.AuthPlugin, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	// Create a new AWS session
	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use env config
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	client := cognitoidentityprovider.New(sess)

	return &CognitoPlugin{
		client: client,
	}, nil
}
