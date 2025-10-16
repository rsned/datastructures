package tree

import (
	"bytes"

	"golang.org/x/exp/constraints"
)

const treeTypeAVL = "AVL"

// AVL (named after inventors Adelson-Velsky and Landis) is a self-balancing
// binary search tree. In an AVL tree, the heights of the two child subtrees of
// any node differ by at most one. If at any time they differ by more than one,
// rebalancing is performed through rotations to restore this property. This
// ensures that search, insert, and delete operations have a worst-case time
// complexity of O(log n).
//
// This implementation satisfies the Tree[T] interface.
type AVL[T constraints.Ordered] struct {
	root *avlNode[T]
}

// NewAVL creates and returns an empty AVL tree, ready for use.
func NewAVL[T constraints.Ordered]() Tree[T] {
	return &AVL[T]{
		root: nil,
	}
}

// Root returns the root node of the tree, which implements the BinaryTree[T]
// interface. If the tree is empty, it returns nil.
func (t *AVL[T]) Root() BinaryTree[T] {
	if t.root == nil {
		return nil
	}
	return t.root
}

// Insert adds a new value to the AVL tree.
// After insertion, it traverses back up the tree to the root, updating node
// heights and performing rotations as necessary to maintain the AVL balance
// property. Duplicate values are not inserted.
// It returns true if the insertion was successful, or false if the value
// already exists.
func (t *AVL[T]) Insert(v T) bool {
	var inserted bool
	t.root, inserted = t.root.insertInternal(v)
	// The root node should not have a parent.
	if t.root != nil {
		t.root.parent = nil
	}
	return inserted
}

// Delete removes the specified value from the tree.
// Note: This functionality is not yet implemented for the AVL tree.
// It currently returns false for all inputs.
func (t *AVL[T]) Delete(_ T) bool {
	// TODO: Implement AVL deletion.
	return false
}

// Search checks if the specified value exists in the tree.
// It returns true if the value is found, and false otherwise.
func (t *AVL[T]) Search(v T) bool {
	if t == nil || t.root == nil {
		return false
	}
	return t.root.Search(v)
}

// Traverse traverses the tree and returns a read-only channel of values.
// The traversal order is determined by the `tOrder` parameter.
// The channel is closed after all values have been sent.
func (t *AVL[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if t.root != nil {
			traverseBinaryTree(t.root, tOrder, ch)
		}
	}()
	return ch
}

// Height returns the height of the tree. An empty tree has a height of 0,
// and a tree with one node has a height of 1.
func (t *AVL[T]) Height() int {
	if t == nil || t.root == nil {
		return 0
	}
	return t.root.Height()
}

// Clone creates and returns a deep copy of the AVL tree. The new tree is
// completely independent of the original.
func (t *AVL[T]) Clone() Tree[T] {
	if t == nil || t.root == nil {
		return NewAVL[T]()
	}

	clone := &AVL[T]{}
	if clonedRoot, ok := t.root.Clone().(*avlNode[T]); ok {
		clone.root = clonedRoot
	}
	return clone
}

// toTestString is an internal helper function used for testing. It generates a
// string representation of the tree structure that can be used in test cases.
func (t *AVL[T]) toTestString() string {
	var buf bytes.Buffer

	buf.WriteString("tree := &AVL[T]{\n")
	buf.WriteString("\troot: &avlNode[T]{\n")

	t.root.toTestString(&buf, 2)

	buf.WriteString("\t},\n")
	buf.WriteString("}\n")

	return buf.String()
}