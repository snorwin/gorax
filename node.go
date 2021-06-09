package gorax

import "sort"

type node struct {
	key      string
	children []*node
	value    interface{}
}

func (n node) isCompressed() bool {
	return len(n.key) != len(n.children)
}

func (n node) isKey() bool {
	return n.value != nil
}

func (n node) isLeaf() bool {
	return len(n.key) == 0
}

func (n *node) getValue() interface{} {
	if n.isKey() {
		if _, isNil := n.value.(Nil); isNil {
			return nil
		}
	}

	return n.value
}

func (n *node) getKeysWithPrefix(prefix string) []string {
	if n.isCompressed() {
		return []string{prefix + n.key}
	} else {
		ret := make([]string, len(n.key))
		for i, key := range n.key {
			ret[i] = prefix + string(key)
		}

		return ret
	}
}

func (n *node) getChildren() []*node {
	return n.children
}

func (n *node) addChild(key string, child *node) {
	idx := sort.Search(len(n.key), func(i int) bool { return n.key[i] >= key[0] })
	if idx == len(n.key) {
		n.key = n.key + key
		n.children = append(n.children, child)
	} else {
		n.key = n.key[:idx] + key + n.key[idx:]

		n.children = append(n.children[:idx+1], n.children[idx:]...)
		n.children[idx] = child
	}
}

func (n *node) addCompressedChild(key string, child *node) {
	n.key = key
	n.children = []*node{child}
}
