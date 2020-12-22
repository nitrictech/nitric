package dynamodb_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDynamoDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DynamoDb Suite")
}
