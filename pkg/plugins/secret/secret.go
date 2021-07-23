package secret

import "fmt"

type Secret struct {
	Name    string
	Version string
	Value   []byte
}

type SecretService interface {
	Put(secret *Secret) (*SecretPutResponse, error)
	Get(id string, versionId string) (*Secret, error)
}

type SecretPutResponse struct {
	Name    string
	Version string
}

type UnimplementedSecretPlugin struct {
	SecretService
}

var _ SecretService = (*UnimplementedSecretPlugin)(nil)

func (*UnimplementedSecretPlugin) Put(secret *Secret) (*SecretPutResponse, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedSecretPlugin) Get(id string, versionId string) (*Secret, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
