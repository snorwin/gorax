package gorax

type node struct {
	key      []byte
	children []*node
	value    interface{}
}

func (r node) isCompressed() bool {
	return len(r.key) != len(r.children)
}

func (r node) isKey() bool {
	return r.value != nil
}

func (r node) getValue() interface{} {
	return r.value
}

func (r node) getKeysWithPrefix(prefix []byte) [][]byte {
	if r.isCompressed() {
		return [][]byte{append(prefix, r.key...)}
	} else {
		ret := make([][]byte, len(r.key))
		for i, key := range r.key {
			ret[i] = append(ret[i], prefix...)
			ret[i] = append(ret[i], key)
		}
		return ret
	}
}

func (r node) getChildren() []*node {
	return r.children
}
