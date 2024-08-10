package tree

import "golang.org/x/exp/constraints"

// RedBlack Tree.
type RedBlack[T constraints.Ordered] struct {
	root *redBlackNode[T]
}

// NewRedBlack returns an empty Red-Black tree ready to use.
func NewRedBlack[T constraints.Ordered]() Tree[T] {
	return &RedBlack[T]{}
}

// Root returns the root node of the tree.
func (t *RedBlack[T]) Root() BinaryTree[T] {
	return t.root
}

// Insert inserts the node into the tree, growing as needed.
func (t *RedBlack[T]) Insert(v T) bool {
	if t.root == nil {
		t.root = &redBlackNode[T]{
			value: v,
		}
		return true
	}
	return t.root.Insert(v)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
//
// The trees internal structure may be updated.
func (t *RedBlack[T]) Delete(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Delete(v)
}

// Search reports if the given value is in the tree.
func (t *RedBlack[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}

	return t.root.Search(v)
}

// Traverse traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *RedBlack[T]) Traverse(w TraverseOrder) <-chan T {
	return make(chan T)
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *RedBlack[T]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.Height()
}
