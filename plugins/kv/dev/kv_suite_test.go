package kv_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDocuments(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KV Suite")
}
