package gorax

import (
	"sort"
	"strings"
)

// Tree implements a radix tree.
type Tree struct {
	root node
	size int
}

// New returns an empty Tree.
func New() *Tree {
	return &Tree{}
}

// FromMap returns a new Tree containing the keys from an existing map.
func FromMap(values map[string]interface{}) *Tree {
	t := New()
	for k, v := range values {
		t.Insert(k, v)
	}

	return t
}

// ToMap walks the Tree and converts it into a map.
func (t *Tree) ToMap() map[string]interface{} {
	ret := map[string]interface{}{}
	t.Walk(func(key string, value interface{}) bool {
		ret[key] = value

		return false
	})

	return ret
}

// Len returns the number of elements in the Tree.
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

// Get is used to lookup a specific key and returns the value and if it was found.
func (t *Tree) Get(key string) (interface{}, bool) {
	current, idx, split := t.find(key, nil)
	if idx != len(key) || (current.isCompressed() && split != 0) || !current.isKey() {
		return nil, false
	}

	return current.getValue(), true
}

// LongestPrefix is like Get, but instead of an exact match, it will return the longest prefix match.
func (t *Tree) LongestPrefix(prefix string) (string, interface{}, bool) {
	var current *node
	var currentKey string
	t.find(prefix, func(key string, node *node) bool {
		if node.isKey() {
			current = node
			currentKey = key
		}

		return false
	})

	if current == nil {
		return "", nil, false
	}

	return currentKey, current.getValue(), true
}

// Delete deletes a key and returns the previous value and if it was deleted.
func (t *Tree) Delete(key string) (interface{}, bool) {
	var nodes []*node
	current, idx, split := t.find(key, func(_ string, n *node) bool {
		nodes = append(nodes, n)
		return false
	})
	if idx != len(key) || (current.isCompressed() && split != 0) || !current.isKey() {
		return nil, false
	}

	value := current.getValue()
	current.value = nil

	t.size -= 1

	t.delete(nodes[:len(nodes)-1], current)

	return value, true
}

// DeletePrefix deletes the subtree under a prefix Returns how many nodes were deleted.
// Use this to delete large subtrees efficiently.
func (t *Tree) DeletePrefix(prefix string) int {
	var counter int

	var nodes []*node
	current, idx, split := t.find(prefix, func(_ string, n *node) bool {
		nodes = append(nodes, n)
		return false
	})

	if len(prefix) == idx+split {
		walk(current, func(key string, node *node) bool {
			if node.isKey() {
				counter += 1
			}
			return false
		})

		t.size -= counter

		current.key = ""
		current.children = nil

		t.delete(nodes, current)
	}

	return counter
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

// WalkPrefix walks the Tree under a prefix.
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

// WalkPath is used to walk the Tree, but only visiting nodes from the root down to a given leaf.
func (t *Tree) WalkPath(path string, fn WalkFn) {
	t.find(path, func(key string, node *node) bool {
		// call WalkFn
		if node.isKey() {
			return fn(key, node.getValue())
		}

		return false
	})
}

// Minimum returns the minimum value in the Tree.
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

	return string(ret), current.getValue(), current.isKey()
}

// Maximum returns the maximum value in the Tree.
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

func (t *Tree) find(key string, fn func(string, *node) bool) (*node, int, int) {
	current := &t.root

	var idx int
	for len(current.key) > 0 && idx < len(key) {
		if fn != nil {
			// call function if defined
			if fn(key[:idx], current) {
				break
			}
		}

		if current.isCompressed() {
			// match as many chars as possible from the compressed key with the lookup key
			if !strings.HasPrefix(key[idx:], current.key) {
				i := sort.Search(len(current.key), func(i int) bool {
					return !strings.HasPrefix(key[idx:], current.key[:i])
				}) - 1

				return current, idx + i, i
			}

			idx += len(current.key)
			current = current.children[0]
		} else {
			// find a child whose key is matching with the lookup key
			i := sort.Search(len(current.key), func(i int) bool {
				return current.key[i] >= key[idx]
			})
			if i == len(current.key) || current.key[i] != key[idx] {
				// no matching child found - break
				return current, idx, 0
			}

			idx += 1
			current = current.children[i]
		}
	}

	if current.isLeaf() || len(key) == idx {
		if fn != nil {
			// call function if defined
			fn(key[:idx], current)
		}
	}

	return current, idx, 0
}

func (t *Tree) delete(nodes []*node, current *node) {
	var trycompress bool
	if len(current.children) == 0 {
		var child *node
		for current != &t.root {
			child = current

			current = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]

			if current.isKey() || (!current.isCompressed() && len(current.children) != 1) {
				break
			}
		}
		if child != nil {
			current.removeChild(child)
		}

		if len(current.children) == 1 && !current.isKey() {
			trycompress = true
		}
	} else if len(current.children) == 1 {
		trycompress = true
	}

	if trycompress {
		var parent *node
		for {
			if len(nodes) == 0 {
				parent = nil
				break
			}
			parent = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
			if parent.isKey() || (!parent.isCompressed() && len(parent.children) != 1) {
				break
			}
			current = parent
		}

		start := current

		newChild := node{}
		for len(current.children) != 0 {
			newChild.key += current.key
			newChild.children = current.children
			current = current.children[len(current.children)-1]
			if current.isKey() || (!current.isCompressed() && len(current.children) != 1) {
				break
			}
		}
		if newChild.key != "" {
			if parent != nil {
				for i, child := range parent.children {
					if child == start {
						parent.children[i] = &newChild
					}
				}
			} else {
				t.root = newChild
			}
		}
	}
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
