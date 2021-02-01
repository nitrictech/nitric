package mocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type MockCognitoIdentityProvider struct {
	// We only want to implement the methods we actually use...
	cognitoidentityprovideriface.CognitoIdentityProviderAPI

	pools []*MockUserPool
}

func (m *MockCognitoIdentityProvider) Clear() {
	m.pools = make([]*MockUserPool, 0)
}

func (m *MockCognitoIdentityProvider) GetUserPool(name string) *MockUserPool {
	for _, p := range m.pools {
		if name == *p.Pool.Name {
			return p
		}
	}

	return nil
}

func (m *MockCognitoIdentityProvider) SetPools(pools []*MockUserPool) {
	m.pools = pools
}

func (m *MockCognitoIdentityProvider) GetMockUser(poolName, username string) *cognitoidentityprovider.UserType {
	for _, p := range m.pools {
		if *p.Pool.Name == poolName {
			for _, u := range p.Users {
				if *u.Username == username {
					return u
				}
			}
		}
	}

	return nil
}

func (m *MockCognitoIdentityProvider) SignUp(in *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	// Create a new user...
	var userPool *MockUserPool
	// var userPoolClient *cognitoidentityprovider.UserPoolClientType
	for _, p := range m.pools {
		for _, c := range p.Clients {
			if c.ClientId == in.ClientId {
				userPool = p
				// userPoolClient = c
				break
			}
		}
	}

	if userPool == nil {
		return nil, fmt.Errorf("UserPoolClient %s does not exist", *in.ClientId)
	}

	for _, u := range userPool.Users {
		if *in.Username == *u.Username {
			return nil, fmt.Errorf("User with %s username already exists", *in.Username)
		}
	}

	// Now we can create the user
	userPool.Users = append(userPool.Users, &cognitoidentityprovider.UserType{
		Username:   in.Username,
		Attributes: in.UserAttributes,
	})

	return &cognitoidentityprovider.SignUpOutput{
		UserSub:       in.Username,
		UserConfirmed: aws.Bool(false),
	}, nil
}

// TODO: May have to update this for
// different attribute types, for now we're only using email
// so this should be fine...
func (m *MockCognitoIdentityProvider) filterUser(u *cognitoidentityprovider.UserType, filter *string) bool {
	filters := make([]string, 0)
	if filter != nil {
		filters = strings.Split(*filter, "=")
	}

	attribute := strings.Trim(filters[0], " ")
	value := strings.Trim(filters[len(filters)-1], " \"")
	var compareFunc func(string, string) bool
	if len(filters) > 2 {
		// We have a prefix filter
		compareFunc = strings.HasPrefix
	} else if len(filters) > 1 {
		// We have an equality filter
		compareFunc = func(a string, b string) bool {
			return a == b
		}
	}

	if compareFunc != nil {

		for _, att := range u.Attributes {
			if *att.Name == attribute {
				return compareFunc(*att.Value, value)
			}
		}
		// no matching attribute
		return false
	}

	return true
}

func (m *MockCognitoIdentityProvider) ListUsers(in *cognitoidentityprovider.ListUsersInput) (*cognitoidentityprovider.ListUsersOutput, error) {
	for _, p := range m.pools {
		if *in.UserPoolId == *p.Pool.Id {
			users := make([]*cognitoidentityprovider.UserType, 0)

			for _, u := range p.Users {
				// process the filters here as well...
				if m.filterUser(u, in.Filter) {
					users = append(users, u)
				}
			}

			return &cognitoidentityprovider.ListUsersOutput{
				Users: users,
			}, nil
		}
	}

	return nil, fmt.Errorf("UserPool %s does not exist", *in.UserPoolId)
}

func (m *MockCognitoIdentityProvider) ListUserPools(in *cognitoidentityprovider.ListUserPoolsInput) (*cognitoidentityprovider.ListUserPoolsOutput, error) {
	upDescs := make([]*cognitoidentityprovider.UserPoolDescriptionType, 0)
	for _, p := range m.pools {
		upDescs = append(upDescs, &cognitoidentityprovider.UserPoolDescriptionType{
			Id:   p.Pool.Id,
			Name: p.Pool.Name,
		})
	}

	return &cognitoidentityprovider.ListUserPoolsOutput{
		UserPools: upDescs,
	}, nil
}

func (m *MockCognitoIdentityProvider) ListUserPoolClients(in *cognitoidentityprovider.ListUserPoolClientsInput) (*cognitoidentityprovider.ListUserPoolClientsOutput, error) {

	for _, p := range m.pools {
		if *p.Pool.Id == *in.UserPoolId {
			clients := make([]*cognitoidentityprovider.UserPoolClientDescription, 0)
			for _, c := range p.Clients {
				clients = append(clients, &cognitoidentityprovider.UserPoolClientDescription{
					UserPoolId: p.Pool.Id,
					ClientId:   c.ClientId,
					ClientName: c.ClientName,
				})
			}

			return &cognitoidentityprovider.ListUserPoolClientsOutput{
				UserPoolClients: clients,
			}, nil

		}
	}

	return nil, fmt.Errorf("UserPool %s does not exist", *in.UserPoolId)

}

func (m *MockCognitoIdentityProvider) CreateUserPool(in *cognitoidentityprovider.CreateUserPoolInput) (*cognitoidentityprovider.CreateUserPoolOutput, error) {
	// We don't need to See if the pool already exists, cognito will create a new provider with the same name but differentID in this case

	poolId := aws.String(strconv.Itoa(len(m.pools)))

	m.pools = append(m.pools, &MockUserPool{
		Pool: &cognitoidentityprovider.UserPoolType{
			Id:   poolId,
			Name: in.PoolName,
		},
		Clients: make([]*cognitoidentityprovider.UserPoolClientType, 0),
		Users:   make([]*cognitoidentityprovider.UserType, 0),
	})

	return &cognitoidentityprovider.CreateUserPoolOutput{
		UserPool: &cognitoidentityprovider.UserPoolType{
			Id:   poolId,
			Name: in.PoolName,
		},
	}, nil
}

func (m *MockCognitoIdentityProvider) CreateUserPoolClient(in *cognitoidentityprovider.CreateUserPoolClientInput) (*cognitoidentityprovider.CreateUserPoolClientOutput, error) {
	// Find the Userpool for the given poolId
	for _, p := range m.pools {
		if *p.Pool.Id == *in.UserPoolId {
			newUserPoolClient := &cognitoidentityprovider.UserPoolClientType{
				ClientId:   in.ClientName,
				ClientName: in.ClientName,
			}

			// Add it to our list...
			p.Clients = append(p.Clients, newUserPoolClient)

			return &cognitoidentityprovider.CreateUserPoolClientOutput{
				UserPoolClient: newUserPoolClient,
			}, nil
		}
	}

	// TODO: Use actual AWS error in this case...
	return nil, fmt.Errorf("UserPool %s does not exist", *in.UserPoolId)
}

type MockUserPool struct {
	Pool    *cognitoidentityprovider.UserPoolType
	Clients []*cognitoidentityprovider.UserPoolClientType
	Users   []*cognitoidentityprovider.UserType
}

func NewMockCognitoIdentityProvider() *MockCognitoIdentityProvider {
	return &MockCognitoIdentityProvider{
		pools: make([]*MockUserPool, 0),
	}
}
