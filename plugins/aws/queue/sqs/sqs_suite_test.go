package sqs_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSqs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqs Suite")
}
