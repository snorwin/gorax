package gorax

type Tree struct {
	head node
	size int
}

func New() *Tree {
	return &Tree{}
}

func FromMap(values map[string]interface{}) *Tree {
	t := New()
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
	t.walk(func(key []byte, node *node) {
		// call WalkFunc
		if node.isKey() {
			f(string(key), node.getValue())
		}

	})
}

func (t *Tree) Minimum() string {
	current := &t.head

	var ret []byte
	for len(current.key) > 0 {
		if current.isKey() {
			break
		}
		if current.isCompressed() {
			ret = append(ret, current.key...)
		} else {
			ret = append(ret, current.key[0])
		}

		current = current.children[0]
	}

	return string(ret)
}

func (t *Tree) Maximum() string {
	current := &t.head

	var ret []byte
	for len(current.key) > 0 {
		if current.isCompressed() {
			ret = append(ret, current.key...)
			current = current.children[0]
		} else {
			ret = append(ret, current.key[len(current.key)-1])
			current = current.children[len(current.key)-1]
		}
	}

	return string(ret)
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
			newChild := &node{}

			if split == 0 {
				current.children = []*node{
					{
						key:      append([]byte{}, current.key[1:]...),
						children: current.children,
					},
				}

				current.key = []byte{current.key[0]}
				current.addChild(key[idx], newChild)
			} else {
				var oldChild *node
				if len(current.key) == split+1 {
					oldChild = current.children[0]
				} else {
					oldChild = &node{
						key:      append([]byte{}, current.key[split+1:]...),
						children: current.children,
					}
				}

				splitNode := &node{}
				splitNode.addChild(current.key[split], oldChild)
				splitNode.addChild(key[idx], newChild)

				current.key = append([]byte{}, current.key[0:split]...)
				current.children = []*node{splitNode}
			}

			current = newChild
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

		child := &node{}

		// if there are more than one char left and the current key is empty turn it into a compressed node
		if len(current.key) == 0 && len(key) > 1 {
			size = len(key) - idx
			if size > MaxNodeKeySize {
				size = MaxNodeKeySize
			}

			current.addCompressedChild(key[idx:idx+size], child)
		} else {
			size = 1

			current.addChild(key[idx], child)
		}

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

func (t *Tree) walk(fn func([]byte, *node)) {
	nodes := []*node{&t.head}
	keys := [][]byte{[]byte("")}

	for len(nodes) > 0 {
		// pop node
		current := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		// pop key
		key := keys[len(keys)-1]
		keys = keys[:len(keys)-1]

		// call function
		fn(key, current)

		// push child nodes
		nodes = append(nodes, current.getChildren()...)

		// push child keys with current key as prefix
		keys = append(keys, current.getKeysWithPrefix(key)...)
	}
}
