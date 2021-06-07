package gorax

type Tree struct {
	head node
	size int
}

func FromMap(values map[string]interface{}) *Tree {
	t := &Tree{}
	for k, v := range values {
		t.Insert(k, v)
	}
	return t
}

func (t *Tree) ToMap() map[string]interface{} {
	ret := map[string]interface{}{}
	t.Walk(func(key string, value interface{}) {
		ret[key] = value
	})
	return ret
}

func (t *Tree) Len() int {
	return t.size
}

func (t *Tree) Insert(key string, value interface{}) bool {
	if value == nil {
		value = Nil{}
	}

	ok := t.insert([]byte(key), value, true)
	if ok {
		t.size += 1
	}
	return ok
}

func (t *Tree) Get(key string) (interface{}, bool) {
	value, ok := t.get([]byte(key))

	if ok {
		if _, isNil := value.(Nil); isNil {
			return nil, true
		}
	}

	return value, ok
}

type WalkFunc func(key string, value interface{})

func (t *Tree) Walk(f WalkFunc) {
	nodes := []*node{&t.head}
	keys := [][]byte{[]byte("")}

	for len(nodes) > 0 {
		// pop node
		current := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		// pop key
		key := keys[len(keys)-1]
		keys = keys[:len(keys)-1]

		// call WalkFunc
		if current.isKey() {
			f(string(key), current.getValue())
		}

		// push child nodes
		nodes = append(nodes, current.getChildren()...)

		// push child keys with current key as prefix
		keys = append(keys, current.getKeysWithPrefix(key)...)
	}
}

func (t *Tree) insert(key []byte, value interface{}, overwrite bool) bool {
	// find the radix tree as far as possible
	current, idx, split := t.find(key)

	// insert value if key is already part of the tree and not in the middle of a compressed node
	if idx == len(key) && (!current.isCompressed() || split == 0) {
		// update the existing key if there is already one
		if current.isKey() {
			if overwrite {
				current.value = value
			}
			return false
		}

		// insert value
		current.value = value
		return true
	}

	// split compressed node
	if current.isCompressed() {
		if idx != len(key) {
			if split == 0 {
				rightChild := &node{}

				leftChild := &node{
					key:      append([]byte{}, current.key[1:]...),
					children: current.children,
				}

				current.key = []byte{current.key[0], key[idx]}
				current.children = []*node{leftChild, rightChild}

				current = rightChild
			} else {
				rightChild := &node{}

				var leftChild *node
				if len(current.key) == split+1 {
					leftChild = current.children[0]
				} else {
					leftChild = &node{
						key:      append([]byte{}, current.key[split+1:]...),
						children: current.children,
					}
				}

				splitNode := &node{}
				splitNode.key = []byte{current.key[split], key[idx]}
				splitNode.children = []*node{leftChild, rightChild}

				current.key = append([]byte{}, current.key[0:split]...)
				current.children = []*node{splitNode}

				current = rightChild
			}
		} else {
			child := &node{
				key:      append([]byte{}, current.key[split:]...),
				children: current.children,
			}

			current.key = append([]byte{}, current.key[0:split]...)
			current.children = []*node{child}

			current = child
		}
		idx += 1
	}

	// insert missing nodes
	for idx < len(key) {
		var size int

		// if there are more than one char left and the current key is empty turn it into a compressed node
		if len(current.key) == 0 && len(key) > 1 {
			size = len(key) - idx
			if size > MaxNodeKeySize {
				size = MaxNodeKeySize
			}
		} else {
			size = 1
		}

		current.key = append(current.key, key[idx:idx+size]...)

		child := &node{}
		current.children = append(current.children, child)
		current = child

		idx += size
	}

	// insert value
	current.value = value
	return true
}

func (t *Tree) get(key []byte) (interface{}, bool) {
	current, idx, split := t.find(key)
	if idx != len(key) || (current.isCompressed() && split != 0) || !current.isKey() {
		return nil, false
	}

	return current.getValue(), true
}

func (t *Tree) find(key []byte) (*node, int, int) {
	current := &t.head

	var i, j int
	for len(current.key) > 0 && i < len(key) {
		if current.isCompressed() {
			for j = 0; j < len(current.key) && i < len(key); j++ {
				if current.key[j] != key[i] {
					break
				}

				i += 1
			}
			if j != len(current.key) {
				break
			}

			j = 0
		} else {
			for j = 0; j < len(current.key); j++ {
				if current.key[j] == key[i] {
					break
				}
			}
			if j == len(current.key) {
				break
			}

			i += 1
		}

		current = current.children[j]
		j = 0
	}

	return current, i, j
}
