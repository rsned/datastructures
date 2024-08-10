package tree

import "golang.org/x/exp/constraints"

// redBlackNode is the basic node in a Red-Black binary search tree.
type redBlackNode[T constraints.Ordered] struct {
	value T

	isRed bool

	left, right *redBlackNode[T]
}

// HasLeft reports if this node has a Left child.
func (t *redBlackNode[T]) HasLeft() bool {
	return t.left != nil
}

// HasRight reports if this node has a Right child.
func (t *redBlackNode[T]) HasRight() bool {
	return t.right != nil
}

// Left returns this nodes Left child.
func (t *redBlackNode[T]) Left() BinaryTree[T] {
	return t.left
}

// Right returns this nodes Right child.
func (t *redBlackNode[T]) Right() BinaryTree[T] {
	return t.right
}

// Value returns this nodes Value.
func (t *redBlackNode[T]) Value() T {
	return t.value
}

// Metadata returns a string of metadata about this node. In this case if
// the node is a Red or Black node.
func (t *redBlackNode[T]) Metadata() string {
	if t.isRed {
		return "Red"
	}
	return "Black"
}

// Insert inserts the node into the tree, growing as needed, and reports
// if the operation was successful.
func (t *redBlackNode[T]) Insert(v T) bool {
	if t == nil {
		return false
	}

	if v == t.value {
		return false
	}

	if v < t.value {
		if t.left == nil {
			t.left = &redBlackNode[T]{value: v}
			return true
		}
		return t.left.Insert(v)
	}

	if t.right == nil {
		t.right = &redBlackNode[T]{value: v}
		return true
	}
	return t.right.Insert(v)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *redBlackNode[T]) Delete(v T) bool {
	if t == nil {
		return false
	}
	return false
}

// Search reports if the given value is in the tree.
func (t *redBlackNode[T]) Search(v T) bool {
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
func (t *redBlackNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *redBlackNode[T]) Height() int {
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
