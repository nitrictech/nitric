package sdk

import "fmt"

// AuthService - Pure Golang interface
type AuthService interface {
	CreateUser(tenant string, id string, email string, password string) error
}

// UnimplementedAuthPlugin - Unimplemented stub struct for extension for partial implementation of the AuthService
type UnimplementedAuthPlugin struct {
	AuthService
}

// CreateUser - Stub user creation method returning default UNIMPLEMENTED error message
func (s *UnimplementedAuthPlugin) CreateUser(tenant string, id string, email string, password string) error {
	return fmt.Errorf("Unimplemented")
}
