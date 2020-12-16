package membrane_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMembrane(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Membrane Suite")
}
