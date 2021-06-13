package gorax

import (
	"fmt"

	"github.com/emicklei/dot"
)

// ToDOTGraph walks the Tree  and converts it into a dot.Graph
func (t *Tree) ToDOTGraph() *dot.Graph {
	// create new dot graph
	graph := dot.NewGraph(dot.Directed)

	walk(&t.root, func(key string, node *node) bool {
		n := graph.Node(key)
		if node.isKey() {
			// set value in label
			n.Attr("label", fmt.Sprintf("%s|%v", key, node.getValue()))

			// change shape
			n.Attr("shape", "record")
		}

		if node.isLeaf() {
			// leaf nodes are blue
			n.Attr("color", "blue")
		}

		if node.isCompressed() {
			// compressed nodes are green
			n.Attr("color", "green")

			// add compressed edge
			n.Edge(graph.Node(key+node.key)).
				Label(node.key).
				Attr("color", "green")
		} else {
			// add all other edges
			for i := 0; i < len(node.key); i++ {
				n.Edge(graph.Node(key + string(node.key[i]))).
					Label(string(node.key[i]))
			}
		}

		return false
	})

	// create root
	graph.Node("").Attr("shape", "point")

	return graph
}
