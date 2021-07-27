package secret

import "fmt"

type SecretService interface {
	// Put - Creates a new version for a given secret
	Put(*Secret, []byte) (*SecretPutResponse, error)
	// Access - Retrieves the value for a given secret version
	Access(*SecretVersion) (*SecretAccessResponse, error)
}

type UnimplementedSecretPlugin struct {
	SecretService
}

var _ SecretService = (*UnimplementedSecretPlugin)(nil)

func (*UnimplementedSecretPlugin) Put(secret *Secret, value []byte) (*SecretPutResponse, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedSecretPlugin) Access(version *SecretVersion) (*SecretAccessResponse, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
