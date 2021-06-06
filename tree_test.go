package gorax_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/gorax"
)

var _ = Describe("Tree", func() {
	Context("Fuzzy_Insert", func() {
		size := 1000

		t := &gorax.Tree{}
		m := make(map[string]interface{}, size)

		for j := 0; j < size; j++ {
			It("shoud_insert_key_value", func() {
				key := randString(rand.Intn(256))
				value := randString(256)

				_, expected := m[key]
				m[key] = value

				actual := !t.Insert(key, value)
				Ω(actual).Should(Equal(expected))
			})
		}
		AfterEach(func() {
			Ω(t.ToMap()).Should(Equal(m))
		})
	})
	Context("Benchmark_Insert", func() {
		Measure("Benchmark__100000", func(b Benchmarker) {
			size := 100000

			m := make(map[string]interface{}, size)
			for i := 0; i < size; i++ {
				m[randString(rand.Intn(32))] = randString(rand.Intn(32))
			}

			t := gorax.FromMap(m)

			key := randString(rand.Intn(32))
			value := randString(rand.Intn(32))

			runtime := b.Time("runtime", func() {
				t.Insert(key, value)
			})

			actual, ok := t.Get(key)
			Ω(ok).Should(BeTrue())
			Ω(actual).Should(Equal(value))

			b.RecordValueWithPrecision("runtime [μs]", float64(runtime.Microseconds()), "μs", 3)
		}, 100)
	})
})

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
