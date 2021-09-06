package azblob_service

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAzblob(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Azblob Suite")
}
