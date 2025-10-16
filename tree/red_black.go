package tree

import "golang.org/x/exp/constraints"

const treeTypeRedBlack = "Red-Black"

// RedBlack implements a self-balancing binary search tree that uses node
// coloring (red or black) to ensure that the tree remains approximately
// balanced during insertions and deletions. This balancing ensures that
// fundamental operations like search, insert, and delete have a worst-case
// time complexity of O(log n).
//
// The Red-Black tree properties are:
//  1. Every node is either red or black.
//  2. The root is always black.
//  3. There are no two adjacent red nodes (a red node cannot have a red parent
//     or a red child).
//  4. Every path from a node to a nil descendant (leaf) contains the same
//     number of black nodes.
//
// This implementation satisfies the Tree[T] interface.
type RedBlack[T constraints.Ordered] struct {
	root *redBlackNode[T]
}

// NewRedBlack creates and returns an empty Red-Black tree, ready for use.
func NewRedBlack[T constraints.Ordered]() Tree[T] {
	return &RedBlack[T]{
		root: nil,
	}
}

// Root returns the root node of the tree, which implements the BinaryTree[T]
// interface. If the tree is empty, it returns nil.
func (t *RedBlack[T]) Root() BinaryTree[T] {
	if t.root == nil {
		return nil
	}
	return t.root
}

// Insert adds a new value to the Red-Black tree.
// The new node is initially colored red. The tree then performs a series of
// checks and rebalancing operations (recoloring and rotations) to ensure that
// all Red-Black properties are maintained.
// It returns true if the insertion was successful, or false if the value
// already exists.
func (t *RedBlack[T]) Insert(v T) bool {
	if t.root == nil {
		// The root of a Red-Black tree is always black.
		t.root = &redBlackNode[T]{value: v, isRed: false}
		return true
	}

	// The internal insert method handles the recursive insertion and rebalancing.
	// We need to find the root of the tree after potential rotations.
	inserted := t.root.Insert(v)
	for t.root.parent != nil {
		t.root = t.root.parent
	}
	// The root must always be black.
	t.root.isRed = false
	return inserted
}

// Delete removes the specified value from the tree.
// Note: This functionality is not yet fully implemented for the Red-Black tree.
// It may not correctly rebalance the tree after deletion.
func (t *RedBlack[T]) Delete(v T) bool {
	if t.root == nil {
		return false
	}
	// TODO: Implement proper Red-Black deletion and rebalancing.
	return t.root.Delete(v)
}

// Search checks if the specified value exists in the tree.
// It returns true if the value is found, and false otherwise.
func (t *RedBlack[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Search(v)
}

// Traverse traverses the tree and returns a read-only channel of values.
// The traversal order is determined by the `tOrder` parameter.
// The channel is closed after all values have been sent.
func (t *RedBlack[T]) Traverse(tOrder TraverseOrder) <-chan T {
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
func (t *RedBlack[T]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.Height()
}

// Clone creates and returns a deep copy of the Red-Black tree. The new tree
// is completely independent of the original.
func (t *RedBlack[T]) Clone() Tree[T] {
	if t == nil || t.root == nil {
		return NewRedBlack[T]()
	}

	clone := &RedBlack[T]{}
	if clonedRoot, ok := t.root.Clone().(*redBlackNode[T]); ok {
		clone.root = clonedRoot
	}
	return clone
}