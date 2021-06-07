package gorax_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/gorax"
)

var _ = Describe("Tree", func() {
	Context("Insert", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = &gorax.Tree{}
		})
		It("should_insert_nil", func() {
			key := "foo"

			ok := t.Insert(key, nil)
			Ω(ok).Should(BeTrue())

			value, ok := t.Get(key)
			Ω(ok).Should(BeTrue())
			Ω(value).Should(BeNil())
		})
		It("should_overwrite", func() {
			key := "foo"

			Ω(t.Insert(key, "old")).Should(BeTrue())
			Ω(t.Insert(key, "new")).Should(BeFalse())

			value, ok := t.Get(key)
			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal("new"))
			Ω(t.Len()).Should(Equal(1))
		})
	})
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
			Ω(t.Len()).Should(Equal(len(m)))
		})
	})
	Context("Benchmark_Insert", func() {
		size := 100000
		Measure(fmt.Sprintf("Benchmark_Insert_%d", size), func(b Benchmarker) {
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
	Context("Benchmark_Get", func() {
		size := 100000

		m := make(map[string]interface{}, size)
		keys := make([]string, size)
		for i := 0; i < size; i++ {
			keys[i] = randString(rand.Intn(32))
			m[keys[i]] = randString(rand.Intn(32))
		}

		t := gorax.FromMap(m)

		Measure(fmt.Sprintf("Benchmark_Get_%d", size), func(b Benchmarker) {
			key := keys[rand.Intn(len(keys))]

			var value interface{}
			var ok bool
			runtime := b.Time("runtime", func() {
				value, ok = t.Get(key)
			})

			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal(m[key]))

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
