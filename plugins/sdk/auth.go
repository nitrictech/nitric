package sdk

import "fmt"

// UserService - Pure Golang interface
type UserService interface {
	CreateUser(tenant string, id string, email string, password string) error
}

// UnimplementedAuthPlugin - Unimplemented stub struct for extension for partial implementation of the UserService
type UnimplementedAuthPlugin struct {
	UserService
}

// CreateUser - Stub user creation method returning default UNIMPLEMENTED error message
func (s *UnimplementedAuthPlugin) CreateUser(tenant string, id string, email string, password string) error {
	return fmt.Errorf("Unimplemented")
}
