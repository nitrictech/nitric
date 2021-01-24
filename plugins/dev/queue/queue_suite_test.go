package queue_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Queue Suite")
}
