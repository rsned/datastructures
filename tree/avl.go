package tree

import (
	"bytes"

	"golang.org/x/exp/constraints"
)

// AVL tree (named after inventors Adelson-Velsky and Landis) is a
// self-balancing binary search tree. In an AVL tree, the heights of the two
// child subtrees of any node differ by at most one; if at any time they differ
// by more than one, rebalancing is done to restore this property.
type AVL[T constraints.Ordered] struct {
	root *avlNode[T]
}

// NewAVL returns an empty AVL tree ready to use.
func NewAVL[T constraints.Ordered]() Tree[T] {
	return &AVL[T]{
		root: nil,
	}
}

// Root returns the root node of the tree.
func (t *AVL[T]) Root() BinaryTree[T] {
	return t.root
}

// Insert inserts the node into the tree, growing as needed.
func (t *AVL[T]) Insert(v T) bool {
	var inserted bool
	t.root, inserted = t.root.insertInternal(v)
	// Ensure root has no parent
	if t.root != nil {
		t.root.parent = nil
	}

	return inserted
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *AVL[T]) Delete(_ T) bool {
	return false
}

// Search reports if the given value is in the tree.
func (t *AVL[T]) Search(v T) bool {
	if t == nil {
		return false
	}

	return t.root.Search(v)
}

// Traverse traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *AVL[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t.root, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *AVL[T]) Height() int {
	if t == nil {
		return 0
	}

	return t.root.Height()
}

// Clone creates a deep copy of the AVL tree.
func (t *AVL[T]) Clone() Tree[T] {
	if t == nil || t.root == nil {
		return NewAVL[T]()
	}

	clone := &AVL[T]{
		root: nil,
	}
	if clonedRoot, ok := t.root.Clone().(*avlNode[T]); ok {
		clone.root = clonedRoot
	}

	return clone
}

// toTestString prints out this tree with all its properties and children
// ready to copy and paste into test code.
// NOTE: This does not determine the exact type of T this instance is. It
// simply prints the types as [T]. Updating is left to the consumer.
func (t *AVL[T]) toTestString() string {
	var buf bytes.Buffer

	buf.WriteString("tree := &AVL[T]{\n")
	buf.WriteString("\troot: &avlNode[T]{\n")

	t.root.toTestString(&buf, 2)

	buf.WriteString("\t},\n")
	buf.WriteString("}\n")

	return buf.String()
}
