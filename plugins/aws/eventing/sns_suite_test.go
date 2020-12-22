package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSns(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sns Suite")
}
