package pubsub_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPubsub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pubsub Suite")
}
