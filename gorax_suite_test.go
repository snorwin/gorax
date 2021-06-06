package gorax_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGorax(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorax Test Suite")
}
