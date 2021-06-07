package gorax_test

import (
	_ "embed"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/gorax"
)

//go:embed example.dot
var example string

var _ = Describe("Tree", func() {
	Context("ToDotGraph", func() {
		It("should_generate_graph_regression_test", func() {
			t := gorax.FromMap(map[string]interface{}{
				"alligator":     nil,
				"alien":         1,
				"baloon":        2,
				"chromodynamic": 3,
				"romane":        4,
				"romanus":       5,
				"romulus":       6,
				"rubens":        7,
				"ruber":         8,
				"rubicon":       9,
				"rubicundus":    "a",
				"all":           "b",
				"rub":           "c",
				"ba":            "d",
			})

			Ω(strings.ReplaceAll(t.ToDotGraph().String(), "\t\n", "\n")).Should(Equal(example))
		})
	})
})
