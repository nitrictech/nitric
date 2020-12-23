package firestore_plugin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFirestore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firestore Suite")
}
