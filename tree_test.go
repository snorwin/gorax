package gorax_test

import (
	"math/rand"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/gorax"
)

const (
	LetterBytes      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	LetterBytesSmall = "abcdef"

	TestKeyLength   = 256
	TestValueLength = 256

	FuzzyTestSize = 1000

	BenchmarkSamples  = 10
	BenchmarkTreeSize = 100000
)

var _ = Describe("Tree", func() {
	Context("Len", func() {
		It("should_not_fail_if_empty", func() {
			Ω(gorax.New().Len()).Should(Equal(0))
		})
	})
	Context("Insert", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
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
		It("should_insert_key_value_fuzzy", func() {
			size := FuzzyTestSize

			m := make(map[string]interface{}, size)
			for j := 0; j < size; j++ {
				key := randString(rand.Intn(TestKeyLength), LetterBytes)
				value := randString(TestValueLength, LetterBytes)

				_, expected := m[key]
				m[key] = value

				actual := !t.Insert(key, value)
				Ω(actual).Should(Equal(expected))

				Ω(t.ToMap()).Should(Equal(m))
				Ω(t.Len()).Should(Equal(len(m)))
			}
		})
		It("should_insert_limited_bytes_key_value_fuzzy", func() {
			size := FuzzyTestSize

			m := make(map[string]interface{}, size)
			for j := 0; j < size; j++ {
				key := randString(rand.Intn(TestKeyLength), LetterBytesSmall)
				value := randString(TestValueLength, LetterBytes)

				_, expected := m[key]
				m[key] = value

				actual := !t.Insert(key, value)
				Ω(actual).Should(Equal(expected))

				Ω(t.ToMap()).Should(Equal(m))
				Ω(t.Len()).Should(Equal(len(m)))
			}
		})
		Measure("Benchmark", func(b Benchmarker) {
			size := BenchmarkTreeSize

			m := make(map[string]interface{}, size)
			for i := 0; i < size; i++ {
				m[randString(rand.Intn(TestKeyLength), LetterBytes)] = randString(rand.Intn(32), LetterBytes)
			}

			t = gorax.FromMap(m)

			key := randString(rand.Intn(TestKeyLength), LetterBytes)
			value := randString(rand.Intn(TestValueLength), LetterBytes)

			runtime := b.Time("runtime", func() {
				t.Insert(key, value)
			})

			actual, ok := t.Get(key)
			Ω(ok).Should(BeTrue())
			Ω(actual).Should(Equal(value))

			b.RecordValueWithPrecision("runtime [μs]", float64(runtime.Microseconds()), "μs", 3)
		}, BenchmarkSamples)
	})
	Context("Get", func() {
		It("should_not_fail_if_empty", func() {
			value, ok := gorax.New().Get("foo")
			Ω(ok).Should(BeFalse())
			Ω(value).Should(BeNil())
		})
		Measure("Benchmark", func(b Benchmarker) {
			size := BenchmarkTreeSize

			m := make(map[string]interface{}, size)
			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength), LetterBytes)
				m[keys[i]] = randString(rand.Intn(TestValueLength), LetterBytes)
			}

			t := gorax.FromMap(m)

			key := keys[rand.Intn(len(keys))]

			var value interface{}
			var ok bool
			runtime := b.Time("runtime", func() {
				value, ok = t.Get(key)
			})

			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal(m[key]))

			b.RecordValueWithPrecision("runtime [μs]", float64(runtime.Microseconds()), "μs", 3)
		}, BenchmarkSamples)
	})
	Context("Minimum", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			Ω(t.Minimum()).Should(Equal(""))
		})
		It("should_find_minimum", func() {
			t = gorax.FromMap(map[string]interface{}{
				"foo":       1,
				"foobar":    2,
				"foofoo":    3,
				"barbar":    nil,
				"barfoo":    "foo",
				"barbarbar": "bar",
				"foobarfoo": "foo",
			})
			Ω(t.Minimum()).Should(Equal("barbar"))
		})
		It("should_find_minimum_fuzzy", func() {
			size := FuzzyTestSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength), LetterBytes)
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)
			Ω(t.Minimum()).Should(Equal(keys[0]))
		})
		Measure("Benchmark", func(b Benchmarker) {
			size := BenchmarkTreeSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength), LetterBytes)
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)

			var key string
			runtime := b.Time("runtime", func() {
				key = t.Minimum()
			})
			Ω(key).Should(Equal(keys[0]))

			b.RecordValueWithPrecision("runtime [μs]", float64(runtime.Microseconds()), "μs", 3)
		}, BenchmarkSamples)
	})
	Context("Maximum", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			Ω(t.Maximum()).Should(Equal(""))
		})
		It("should_find_maximum", func() {
			t = gorax.FromMap(map[string]interface{}{
				"foo":       1,
				"foobar":    2,
				"foofoo":    3,
				"barbar":    nil,
				"barfoo":    "foo",
				"barbarbar": "bar",
				"foobarfoo": "foo",
			})
			Ω(t.Maximum()).Should(Equal("foofoo"))
		})
		It("should_find_maximum_fuzzy", func() {
			size := FuzzyTestSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength), LetterBytes)
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)
			Ω(t.Maximum()).Should(Equal(keys[len(keys)-1]))
		})
		Measure("Benchmark", func(b Benchmarker) {
			size := BenchmarkTreeSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength), LetterBytes)
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)

			var key string
			runtime := b.Time("runtime", func() {
				key = t.Maximum()
			})
			Ω(key).Should(Equal(keys[len(keys)-1]))

			b.RecordValueWithPrecision("runtime [μs]", float64(runtime.Microseconds()), "μs", 3)
		}, BenchmarkSamples)
	})
})

func randString(n int, bytes string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = bytes[rand.Intn(len(bytes))]
	}
	return string(b)
}
