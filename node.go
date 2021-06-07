package gorax

import "sort"

type node struct {
	key      []byte
	children []*node
	value    interface{}
}

func (n node) isCompressed() bool {
	return len(n.key) != len(n.children)
}

func (n node) isKey() bool {
	return n.value != nil
}

func (n *node) getValue() interface{} {
	return n.value
}

func (n *node) getKeysWithPrefix(prefix []byte) [][]byte {
	if n.isCompressed() {
		return [][]byte{append(prefix, n.key...)}
	} else {
		ret := make([][]byte, len(n.key))
		for i, key := range n.key {
			ret[i] = append(ret[i], prefix...)
			ret[i] = append(ret[i], key)
		}
		return ret
	}
}

func (n *node) getChildren() []*node {
	return n.children
}

func (n *node) deleteChildren() {
	n.key = []byte{}
	n.children = []*node{}
}

func (n *node) addChild(key byte, child *node) {
	idx := sort.Search(len(n.key), func(i int) bool { return n.key[i] >= key })
	if idx == len(n.key) {
		n.key = append(n.key, key)
		n.children = append(n.children, child)
	} else {
		n.key = append(n.key[:idx+1], n.key[idx:]...)
		n.key[idx] = key

		n.children = append(n.children[:idx+1], n.children[idx:]...)
		n.children[idx] = child
	}
}

func (n *node) addCompressedChild(key []byte, child *node) {
	n.key = append(n.key, key...)
	n.children = append(n.children, child)
}
