package cognito_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCognito(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cognito Suite")
}
