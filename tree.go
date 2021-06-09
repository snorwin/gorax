package gorax

import "sort"

// Tree implements a radix tree,
type Tree struct {
	root node
	size int
}

// New returns an empty Tree
func New() *Tree {
	return &Tree{}
}

// FromMap returns a new Tree containing the keys from an existing map
func FromMap(values map[string]interface{}) *Tree {
	t := New()
	for k, v := range values {
		t.Insert(k, v)
	}

	return t
}

// ToMap walks the Tree and converts it into a map
func (t *Tree) ToMap() map[string]interface{} {
	ret := map[string]interface{}{}
	t.Walk(func(key string, value interface{}) bool {
		ret[key] = value

		return false
	})

	return ret
}

// Len returns the number of elements in the Tree
func (t *Tree) Len() int {
	return t.size
}

// Insert adds a new entry or updates an existing entry. Returns 'true' if entry was added.
func (t *Tree) Insert(key string, value interface{}) bool {
	if value == nil {
		value = Nil{}
	}

	ok := t.insert(key, value, true)
	if ok {
		t.size += 1
	}
	return ok
}

// Get is used to lookup a specific key and returns the value and if it was found
func (t *Tree) Get(key string) (interface{}, bool) {
	value, ok := t.get(key)

	if ok {
		if _, isNil := value.(Nil); isNil {
			return nil, true
		}
	}

	return value, ok
}

// LongestPrefix is like Get, but instead of an exact match, it will return the longest prefix match.
func (t *Tree) LongestPrefix(prefix string) (string, interface{}, bool) {
	var current *node
	var currentKey string
	t.find(prefix, func(key string, n *node) bool {
		current = n
		currentKey = key

		return true
	})

	if current == nil {
		return "", nil, false
	}

	if _, isNil := current.getValue().(Nil); isNil {
		return currentKey, nil, true
	}

	return currentKey, current.getValue(), true
}

// Delete deletes a key and returns the previous value and if it was deleted
func (t *Tree) Delete(key string) bool {
	// TODO
	return false
}

// DeletePrefix deletes the subtree under a prefix Returns how many nodes were deleted.
// Use this to delete large subtrees efficiently.
func (t *Tree) DeletePrefix(prefix string) int {
	// TODO
	return 0
}

// WalkFn is used when walking the Tree. Takes a key and value, returning 'true' if iteration should be terminated.
type WalkFn func(key string, value interface{}) bool

// Walk walks the Tree
func (t *Tree) Walk(fn WalkFn) {
	walk(&t.root, func(key string, node *node) bool {
		// call WalkFn
		if node.isKey() {
			return fn(key, node.getValue())
		}

		return false
	})
}

// WalkPrefix walks the tree under a prefix
func (t *Tree) WalkPrefix(prefix string, fn WalkFn) {
	current, idx, split := t.find(prefix, nil)
	if len(prefix) == idx+split {
		walk(current, func(key string, node *node) bool {
			// call WalkFn
			if node.isKey() {
				return fn(prefix+key, node.getValue())
			}

			return false
		})
	}
}

// WalkPath is used to walk the tree, but only visiting nodes from the root down to a given leaf.
func (t *Tree) WalkPath(path string, fn WalkFn) {
	t.find(path, func(key string, node *node) bool {
		// call WalkFn
		if node.isKey() {
			return fn(key, node.getValue())
		}

		return false
	})
}

// Minimum returns the minimum value in the Tree
func (t *Tree) Minimum() (string, interface{}, bool) {
	current := &t.root

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

	if _, isNil := current.getValue().(Nil); isNil {
		return string(ret), nil, true
	}

	return string(ret), current.getValue(), current.isKey()
}

// Maximum returns the maximum value in the Tree
func (t *Tree) Maximum() (string, interface{}, bool) {
	current := &t.root

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

	if _, isNil := current.getValue().(Nil); isNil {
		return string(ret), nil, true
	}

	return string(ret), current.getValue(), current.isKey()
}

func (t *Tree) insert(key string, value interface{}, overwrite bool) bool {
	// find the radix tree as far as possible
	current, idx, split := t.find(key, nil)

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
						key:      current.key[1:],
						children: current.children,
					},
				}

				current.key = string(current.key[0])
				current.addChild(string(key[idx]), newChild)
			} else {
				var oldChild *node
				if len(current.key) == split+1 {
					oldChild = current.children[0]
				} else {
					oldChild = &node{
						key:      current.key[split+1:],
						children: current.children,
					}
				}

				splitNode := &node{}
				splitNode.addChild(string(current.key[split]), oldChild)
				splitNode.addChild(string(key[idx]), newChild)

				current.key = current.key[0:split]
				current.children = []*node{splitNode}
			}

			current = newChild
		} else {
			child := &node{
				key:      current.key[split:],
				children: current.children,
			}

			current.key = current.key[0:split]
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

			current.addCompressedChild(key[idx:idx+size], child)
		} else {
			size = 1

			current.addChild(string(key[idx]), child)
		}

		current = child

		idx += size
	}

	// insert value
	current.value = value
	return true
}

func (t *Tree) get(key string) (interface{}, bool) {
	current, idx, split := t.find(key, nil)
	if idx != len(key) || (current.isCompressed() && split != 0) || !current.isKey() {
		return nil, false
	}

	return current.getValue(), true
}

func (t *Tree) find(key string, fn func(string, *node) bool) (*node, int, int) {
	current := &t.root

	var i, j int
	for len(current.key) > 0 && i < len(key) {
		if fn != nil {
			// call function if defined
			if fn(key[:i], current) {
				break
			}
		}

		if current.isCompressed() {
			// match as many chars as possible from the compressed key with the lookup key
			for j = 0; j < len(current.key) && i < len(key); j++ {
				if current.key[j] != key[i] {
					break
				}

				i += 1
			}

			if j != len(current.key) {
				// not the entire compressed key matched with the lookup key - break
				break
			}

			j = 0
		} else {
			// find a child whose key is matching with the lookup key
			j = sort.Search(len(current.key), func(idx int) bool {
				return current.key[idx] >= key[i]
			})
			if j == len(current.key) || current.key[j] != key[i] {
				// no matching child found - break
				break
			}

			i += 1
		}

		current = current.children[j]
		j = 0
	}

	if current.isLeaf() || (len(key) == i && j == 0) {
		if fn != nil {
			// call function if defined
			fn(key[:i], current)
		}
	}

	return current, i, j
}

func walk(start *node, fn func(string, *node) bool) {
	nodes := []*node{start}
	keys := []string{""}

	for len(nodes) > 0 {
		// pop node
		current := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		// pop key
		key := keys[len(keys)-1]
		keys = keys[:len(keys)-1]

		// call function
		if fn(key, current) {
			break
		}

		// push child nodes
		nodes = append(nodes, current.getChildren()...)

		// push child keys with current key as prefix
		keys = append(keys, current.getKeysWithPrefix(key)...)
	}
}
