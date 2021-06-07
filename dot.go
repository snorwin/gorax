package gorax

import (
	"fmt"

	"github.com/emicklei/dot"
)

func (t *Tree) ToDotGraph() *dot.Graph {
	// create new dot graph
	graph := dot.NewGraph(dot.Directed)

	t.walk(func(key []byte, node *node) {
		n := graph.Node(string(key))
		if node.isKey() {
			// set value in label
			n.Attr("label", fmt.Sprintf("%s|%v", key, node.getValue()))

			// change shape and color
			n.Attr("color", "blue").
				Attr("shape", "record")
		}

		if node.isCompressed() {
			// change color of compressed nodes
			n.Attr("color", "green")

			// add compressed edge
			n.Edge(graph.Node(string(append(key, node.key...)))).
				Label(node.key).
				Attr("color", "green")
		} else {
			// add all edges
			for i := 0; i < len(node.key); i++ {
				n.Edge(graph.Node(string(append(key, node.key[i])))).
					Label(string(node.key[i]))
			}
		}
	})

	// create root
	graph.Node("").Attr("shape", "point")

	return graph
}
