package gorax_test

import (
	"math/rand"
	"sort"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/gorax"
)

const (
	FuzzyTestSize   = 1000
	FuzzyMaxKeySize = 256
)

var _ = Describe("Fuzzy Tests", func() {
	Context("Insert/Get/Delete", func() {
		It("should_insert_get_and_delete", func() {
			t := gorax.New()

			m := make(map[string]interface{}, FuzzyTestSize)
			for i := 0; i < FuzzyTestSize; i++ {
				key := randString(rand.Intn(FuzzyMaxKeySize))
				value := randInteface()

				_, expected := m[key]
				m[key] = value

				actual := !t.Insert(key, value)
				Ω(actual).Should(Equal(expected))

				Ω(t.ToMap()).Should(Equal(m))
				Ω(t.Len()).Should(Equal(len(m)))
			}

			for key, value := range m {
				actual, ok := t.Get(key)
				Ω(ok).Should(BeTrue())

				if value == nil {
					Ω(actual).Should(BeNil())
				} else {
					Ω(actual).Should(Equal(value))
				}
			}

			for key := range m {
				expected := m[key]
				delete(m, key)

				actual, ok := t.Delete(key)
				Ω(ok).Should(BeTrue())

				if expected == nil {
					Ω(actual).Should(BeNil())
				} else {
					Ω(actual).Should(Equal(expected))
				}

				Ω(t.ToMap()).Should(Equal(m))
				Ω(t.Len()).Should(Equal(len(m)))
			}
		})
	})
	Context("Minimum/Maximum", func() {
		var (
			t *gorax.Tree
			m map[string]interface{}

			keys []string
		)
		BeforeEach(func() {
			m = make(map[string]interface{}, FuzzyTestSize)
			for i := 0; i < FuzzyTestSize; i++ {
				m[randString(rand.Intn(FuzzyMaxKeySize))] = randInteface()
			}

			t = gorax.FromMap(m)

			keys = []string{}
			for key := range m {
				keys = append(keys, key)
			}

			sort.Strings(keys)
		})
		It("should_find_minimum", func() {
			for i := 0; i < len(keys); i++ {
				expected := m[keys[i]]

				key, actual, ok := t.Minimum()
				Ω(ok).Should(BeTrue())
				Ω(key).Should(Equal(keys[i]))

				if expected == nil {
					Ω(actual).Should(BeNil())
				} else {
					Ω(actual).Should(Equal(expected))
				}

				t.Delete(key)
			}
		})
		It("should_find_maximum", func() {
			for i := len(keys) - 1; i >= 0; i-- {
				expected := m[keys[i]]

				key, actual, ok := t.Maximum()
				Ω(ok).Should(BeTrue())
				Ω(key).Should(Equal(keys[i]))

				if expected == nil {
					Ω(actual).Should(BeNil())
				} else {
					Ω(actual).Should(Equal(expected))
				}

				t.Delete(key)
			}
		})
	})
	Context("Walk", func() {
		It("should_walk_prefix", func() {
			for i := 0; i < FuzzyTestSize; i++ {
				prefix := randString(rand.Intn(24))

				size := 200
				m := make(map[string]interface{}, 2*size)
				expected := make(map[string]interface{})
				for j := 0; j < size; j++ {
					key := randString(rand.Intn(FuzzyMaxKeySize))
					m[key] = randInteface()

					if strings.HasPrefix(key, prefix) {
						expected[key] = m[key]
					}
				}
				for j := 0; j < size; j++ {
					key := prefix + randString(rand.Intn(FuzzyMaxKeySize))
					m[key] = randInteface()
					expected[key] = m[key]
				}

				t := gorax.FromMap(m)

				actual := make(map[string]interface{})
				t.WalkPrefix(prefix, func(key string, value interface{}) bool {
					actual[key] = value

					return false
				})

				Ω(actual).Should(Equal(expected))
			}
		})
		It("should_walk_path", func() {
			for i := 0; i < FuzzyTestSize; i++ {
				slice := make([]string, 24)
				for j := 0; j < len(slice); j++ {
					slice[j] = randString(rand.Intn(FuzzyMaxKeySize))
				}
				m := map[string]interface{}{}
				for j := len(slice); j >= 0; j-- {
					m[strings.Join(slice[:j], "")] = i
				}

				t := gorax.FromMap(m)

				path := strings.Join(slice, "")
				for j := 0; j <= len(path); j++ {
					expected := map[string]interface{}{}
					for k, v := range m {
						if strings.HasPrefix(path[:j], k) {
							expected[k] = v
						}
					}

					actual := make(map[string]interface{})
					t.WalkPath(path[:j], func(key string, value interface{}) bool {
						actual[key] = value

						return false
					})

					Ω(actual).Should(Equal(expected))
				}
			}
		})
	})
})
