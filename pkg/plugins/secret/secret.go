package secret

import "fmt"

type Secret struct {
	Name   string
	Value  []byte
	Labels map[string]string
}

type SecretService interface {
	Put(secret *Secret) (*SecretPutResponse, error)
	Get(id string, versionId string) (*Secret, error)
}

type SecretPutResponse struct {
	Id        string
	VersionId string
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
