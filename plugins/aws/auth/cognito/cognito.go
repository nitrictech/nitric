package cognito_plugin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// CognitoPlugin - Cognito implementation of the Nitric Auth plugin interface
type CognitoPlugin struct {
	sdk.UnimplementedAuthPlugin
	client *cognitoidentityprovider.CognitoIdentityProvider
}

const DEFAULT_USER_POOL_CLIENT_NAME = "Nitric"

// Get the client id for a given user pool
func (s *CognitoPlugin) findOrCreateUserPoolForTenant(tenant string) (*string, error) {
	// TODO: Need to list over UserPools first, and then use the default NitricClient from that pool

	out, err := s.client.ListUserPools(&cognitoidentityprovider.ListUserPoolsInput{
		// FIXME: Need to implement result paging and supporting unlimited number of tenants...
		MaxResults: aws.Int64(60),
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to search existing user pools: %v", err)
	}

	var poolID *string
	fmt.Println(fmt.Sprintf("Enumerating %d user pools", len(out.UserPools)))
	for _, pool := range out.UserPools {
		fmt.Println(fmt.Sprintf("Comparing %s -> %s", *pool.Name, tenant))
		if *pool.Name == tenant {
			poolID = pool.Id
			break
		}
	}

	if poolID == nil {
		// Create a New UserPool for this tenant
		out, err := s.client.CreateUserPool(&cognitoidentityprovider.CreateUserPoolInput{
			PoolName: &tenant,
		})

		if err != nil {
			return nil, fmt.Errorf("Could not create new user pool for tenant: %v", err)
		}

		poolID = out.UserPool.Id
	}

	// Attempt to find the default NitricClient for this tenant
	upOut, err := s.client.ListUserPoolClients(&cognitoidentityprovider.ListUserPoolClientsInput{
		UserPoolId: poolID,
	})

	if err != nil {
		return nil, fmt.Errorf("Error retrieving existing user pool clients: %v", err)
	}
	// Attempt to find the Nitric SDK client for this tenant...
	var pClientID *string
	for _, pClient := range upOut.UserPoolClients {
		if *pClient.ClientName == DEFAULT_USER_POOL_CLIENT_NAME {
			pClientID = pClient.ClientId
		}
	}

	if pClientID == nil {
		upClient, err := s.client.CreateUserPoolClient(&cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId:        poolID,
			ClientName:        aws.String(DEFAULT_USER_POOL_CLIENT_NAME),
			ExplicitAuthFlows: []*string{aws.String("ALLOW_USER_PASSWORD_AUTH"), aws.String("ALLOW_REFRESH_TOKEN_AUTH")},
			GenerateSecret:    aws.Bool(false),
		})

		if err != nil {
			return nil, fmt.Errorf("Error creating new UserPoolClient for Nitric %v", err)
		}

		pClientID = upClient.UserPoolClient.ClientId
	}

	return pClientID, nil
}

// CreateUser - Creates a new user in AWS cognito
func (s *CognitoPlugin) CreateUser(tenant string, id string, email string, password string) error {
	// Attempt to sign up the user...
	upClient, err := s.findOrCreateUserPoolForTenant(tenant)

	if err != nil {
		return fmt.Errorf("Could not SignUp user: %v", err)
	}

	_, err = s.client.SignUp(&cognitoidentityprovider.SignUpInput{
		// TODO: Need to determine the client id in this case
		// For email/password authentication, will likely just do a single user pool for the stack...
		ClientId: upClient,
		Password: &password,
		Username: &id,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			&cognitoidentityprovider.AttributeType{
				Name:  aws.String("email"),
				Value: &email,
			},
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
