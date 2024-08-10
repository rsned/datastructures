package tree

import "golang.org/x/exp/constraints"

// bstNode is the basic node in a binary search tree.
type bstNode[T constraints.Ordered] struct {
	value       T
	left, right *bstNode[T]
}

// HasLeft reports if this node has a Left child.
func (t *bstNode[T]) HasLeft() bool {
	return t.left != nil
}

// HasRight reports if this node has a Right child.
func (t *bstNode[T]) HasRight() bool {
	return t.right != nil
}

// Left returns this nodes Left child.
func (t *bstNode[T]) Left() BinaryTree[T] {
	return t.left
}

// Right returns this nodes Right child.
func (t *bstNode[T]) Right() BinaryTree[T] {
	return t.right
}

// Value returns this nodes Value.
func (t *bstNode[T]) Value() T {
	return t.value
}

// Metadata returns a string of metadata about this node.
// Plain binary search trees have nothing interesting to show.
func (t *bstNode[T]) Metadata() string {
	return ""
}

// Insert inserts the value into the tree, growing as needed, and reports
// if the operation was successful.
func (t *bstNode[T]) Insert(v T) bool {
	if t == nil {
		return false
	}

	// Duplicates are not allowed.
	if v == t.value {
		return false
	}

	if v < t.value {
		if t.left == nil {
			t.left = &bstNode[T]{value: v}
			return true
		}
		return t.left.Insert(v)
	}

	if t.right == nil {
		t.right = &bstNode[T]{value: v}
		return true
	}
	return t.right.Insert(v)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *bstNode[T]) Delete(v T) bool {
	if t == nil {
		return false
	}
	return false
}

// Search reports if the given value is in the tree.
func (t *bstNode[T]) Search(v T) bool {
	if t == nil {
		return false
	}
	if v == t.value {
		return true
	}

	if v < t.value {
		return t.left.Search(v)
	}
	return t.right.Search(v)
}

// Walk traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *bstNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *bstNode[T]) Height() int {
	if t == nil {
		return 0
	}
	lh := t.left.Height()
	rh := t.right.Height()
	if lh > rh {
		return lh + 1
	}
	return rh + 1
}
