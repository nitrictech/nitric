package sdk

import "fmt"

// AuthPlugin - Pure Golang interface
type AuthPlugin interface {
	SignUp(tenant string, id string, email string, password string) error
}

// UnimplementedAuthPlugin - Unimplemented stub struct for extension for partial implementation of the AuthPlugin
type UnimplementedAuthPlugin struct {
	AuthPlugin
}

// CreateUser - Stub user creation method returning default UNIMPLEMENTED error message
func (s *UnimplementedAuthPlugin) SignUp(tenant string, id string, email string, password string) error {
	return fmt.Errorf("Unimplemented")
}
