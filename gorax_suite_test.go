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

func randString(n int) string {
	letters := [][]rune{
		[]rune("abc"),
		[]rune("abcdefghijklmnopqrstuvwxyz"),
		[]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"),
		[]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!@#$%^&*()_+{}:\\|\"<>?/.,';][=-'"),
	}
	runes := letters[rand.Intn(len(letters))]

	r := make([]rune, n)
	for i := range r {
		r[i] = runes[rand.Intn(len(runes))]
	}
	return string(r)
}

func randInteface() interface{} {
	funcs := []func() interface{}{
		func() interface{} {
			return rand.Int()
		},
		func() interface{} {
			return nil
		},
		func() interface{} {
			return []bool{true, false}[rand.Intn(2)]
		},
		func() interface{} {
			return randString(10)
		},
		func() interface{} {
			return struct{ foo string }{"bar"}
		},
	}

	return funcs[rand.Intn(len(funcs))]()
}
