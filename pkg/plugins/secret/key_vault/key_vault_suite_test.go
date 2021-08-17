package key_vault_secret_service

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecretManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Key Vault Suite")
}
