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
	TestKeyLength   = 256
	TestValueLength = 256

	FuzzyTestSize = 1000
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
				key := randString(rand.Intn(TestKeyLength))
				value := randString(TestValueLength)

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
				key := randString(rand.Intn(TestKeyLength))
				value := randString(TestValueLength)

				_, expected := m[key]
				m[key] = value

				actual := !t.Insert(key, value)
				Ω(actual).Should(Equal(expected))

				Ω(t.ToMap()).Should(Equal(m))
				Ω(t.Len()).Should(Equal(len(m)))
			}
		})
	})
	Context("Get", func() {
		It("should_not_fail_if_empty", func() {
			value, ok := gorax.New().Get("foo")
			Ω(ok).Should(BeFalse())
			Ω(value).Should(BeNil())
		})
	})
	Context("WalkPrefix", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			hit := false
			t.WalkPrefix("foo", func(_ string, _ interface{}) bool {
				hit = true

				return false
			})

			Ω(hit).Should(BeFalse())
		})
		It("should_walk_prefix", func() {
			expected := map[string]interface{}{
				"foo":    1,
				"foof":   2,
				"foobar": 3,
				"foofoo": 4,
				"foojin": 6,
			}

			t = gorax.FromMap(expected)
			t.Insert("bar", 7)
			t.Insert("jin", 8)

			actual := make(map[string]interface{})
			t.WalkPrefix("foo", func(key string, value interface{}) bool {
				actual[key] = value

				return false
			})

			Ω(actual).Should(Equal(expected))
		})
		It("should_walk_prefix_and_stop_after_first", func() {
			t = gorax.FromMap(map[string]interface{}{
				"foo":    1,
				"foof":   2,
				"foobar": 3,
				"foofoo": 4,
				"foojin": 6,
				"bar":    7,
				"jin":    8,
			})

			actual := make(map[string]interface{})
			t.WalkPrefix("foo", func(key string, value interface{}) bool {
				actual[key] = value

				return true
			})

			Ω(len(actual)).Should(Equal(1))
			Ω(actual["foo"]).Should(Equal(1))
		})
		It("should_walk_prefix_fuzzy", func() {
			for i := 0; i < FuzzyTestSize; i++ {
				prefix := randString(rand.Intn(24))

				size := 200
				m := make(map[string]interface{}, 2*size)
				expected := make(map[string]interface{})
				for j := 0; j < size; j++ {
					key := randString(rand.Intn(TestKeyLength))
					m[key] = randString(rand.Intn(TestKeyLength))

					if strings.HasPrefix(key, prefix) {
						expected[key] = m[key]
					}
				}
				for j := 0; j < size; j++ {
					key := prefix + randString(rand.Intn(TestKeyLength))
					m[key] = randString(rand.Intn(TestKeyLength))
					expected[key] = m[key]
				}

				t = gorax.FromMap(m)

				actual := make(map[string]interface{})
				t.WalkPrefix(prefix, func(key string, value interface{}) bool {
					actual[key] = value

					return false
				})

				Ω(actual).Should(Equal(expected))
			}
		})
	})
	Context("WalkPath", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_walk_path", func() {
			path := "foo/bar/jin/foofoo/barbar/jinjin"
			slice := strings.Split(path, "/")

			expected := map[string]interface{}{}
			for i := len(slice); i >= 0; i-- {
				key := strings.Join(slice[:i], "/")
				expected[key] = i
				t.Insert(key, i)
			}
			t.Insert("foo/bar/bar", 7)
			t.Insert("foo/jin/foofoo", 8)
			t.Insert("f/j/b", 9)

			actual := make(map[string]interface{})
			t.WalkPath(path, func(key string, value interface{}) bool {
				actual[key] = value

				return false
			})

			Ω(actual).Should(Equal(expected))
		})
		It("should_walk_path_and_stop_after_first", func() {
			path := "foo/bar/jin/foofoo/barbar/jinjin"
			slice := strings.Split(path, "/")

			expected := map[string]interface{}{}
			for i := len(slice); i >= 0; i-- {
				key := strings.Join(slice[:i], "/")
				expected[key] = i
				t.Insert(key, i)
			}
			t.Insert("foo/bar/bar", 7)
			t.Insert("foo/jin/foofoo", 8)
			t.Insert("f/j/b", 9)

			actual := make(map[string]interface{})
			t.WalkPath(path, func(key string, value interface{}) bool {
				actual[key] = value

				return true
			})

			Ω(len(actual)).Should(Equal(1))
			Ω(actual[""]).Should(Equal(0))
		})
		It("should_walk_path_fuzzy", func() {
			for i := 0; i < FuzzyTestSize; i++ {
				slice := make([]string, 24)
				for j := 0; j < len(slice); j++ {
					slice[j] = randString(rand.Intn(TestKeyLength))
				}
				m := map[string]interface{}{}
				for j := len(slice); j >= 0; j-- {
					m[strings.Join(slice[:j], "")] = i
				}

				t = gorax.FromMap(m)

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
	Context("Minimum", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			key, value, ok := t.Minimum()
			Ω(key).Should(Equal(""))
			Ω(value).Should(BeNil())
			Ω(ok).Should(BeFalse())
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
			key, value, ok := t.Minimum()
			Ω(key).Should(Equal("barbar"))
			Ω(value).Should(BeNil())
			Ω(ok).Should(BeTrue())
		})
		It("should_find_minimum_fuzzy", func() {
			size := FuzzyTestSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength))
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)
			key, _, ok := t.Minimum()
			Ω(key).Should(Equal(keys[0]))
			Ω(ok).Should(BeTrue())
		})
	})
	Context("Maximum", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			key, value, ok := t.Maximum()
			Ω(key).Should(Equal(""))
			Ω(value).Should(BeNil())
			Ω(ok).Should(BeFalse())
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
			key, value, ok := t.Maximum()
			Ω(key).Should(Equal("foofoo"))
			Ω(value).Should(Equal(3))
			Ω(ok).Should(BeTrue())
		})
		It("should_find_maximum_fuzzy", func() {
			size := FuzzyTestSize

			keys := make([]string, size)
			for i := 0; i < size; i++ {
				keys[i] = randString(rand.Intn(TestKeyLength))
				t.Insert(keys[i], "")
			}
			sort.Strings(keys)
			key, _, ok := t.Maximum()
			Ω(key).Should(Equal(keys[len(keys)-1]))
			Ω(ok).Should(BeTrue())
		})
	})
	Context("LongestPrefix", func() {
		var (
			t *gorax.Tree
		)
		BeforeEach(func() {
			t = gorax.New()
		})
		It("should_not_fail_if_empty", func() {
			key, value, ok := t.LongestPrefix("foo")
			Ω(key).Should(Equal(""))
			Ω(value).Should(BeNil())
			Ω(ok).Should(BeFalse())
		})
		It("should_find_longest_prefix", func() {
			t = gorax.FromMap(map[string]interface{}{
				"foo":       1,
				"foobar":    2,
				"foofoo":    3,
				"barbar":    nil,
				"barfoo":    "foo",
				"barbarbar": "bar",
				"foobarfoo": "foo",
			})
			key, value, ok := t.LongestPrefix("barbarbarbarbar")
			Ω(key).Should(Equal("barbarbar"))
			Ω(value).Should(Equal("bar"))
			Ω(ok).Should(BeTrue())
		})
	})
})
