package tree

import "golang.org/x/exp/constraints"

// BSTNode is the basic node in a binary search tree.
type BSTNode[T constraints.Ordered] struct {
	value       T
	left, right *BSTNode[T]
}

func (t *BSTNode[T]) HasLeft() bool {
	return t.left != nil
}

func (t *BSTNode[T]) HasRight() bool {
	return t.right != nil
}

func (t *BSTNode[T]) Left() BinaryTree[T] {
	return t.left
}

func (t *BSTNode[T]) Right() BinaryTree[T] {
	return t.right
}

func (t *BSTNode[T]) Value() T {
	return t.value
}

func (t *BSTNode[T]) Metadata() string {
	return ""
}

// Insert inserts the value into the tree, growing as needed, and reports
// if the operation was successful.
func (t *BSTNode[T]) Insert(v T) bool {
	if t == nil {
		return false
	}

	// Duplicates are not allowed.
	if v == t.value {
		return false
	}

	if v < t.value {
		if t.left == nil {
			t.left = &BSTNode[T]{value: v}
			return true
		}
		return t.left.Insert(v)
	}

	if t.right == nil {
		t.right = &BSTNode[T]{value: v}
		return true
	}
	return t.right.Insert(v)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *BSTNode[T]) Delete(v T) bool {
	if t == nil {
		return false
	}
	return false
}

// Search reports if the given value is in the tree.
func (t *BSTNode[T]) Search(v T) bool {
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
func (t *BSTNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *BSTNode[T]) Height() int {
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
