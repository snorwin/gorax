package gorax_test

import (
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGorax(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorax Test Suite")
}

var letterRunes = [][]rune{
	[]rune("abc"),
	[]rune("abcdefghijklmnopqrstuvwxyz"),
	[]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
	[]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!@#$%^&*()_+{}:\\|\"<>?/.,';][=-'"),
}

func randString(n int) string {
	runes := letterRunes[rand.Intn(len(letterRunes))]

	r := make([]rune, n)
	for i := range r {
		r[i] = runes[rand.Intn(len(runes))]
	}
	return string(r)
}
