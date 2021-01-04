package eventing_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEventing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eventing Suite")
}
