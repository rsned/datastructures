package tree

import "golang.org/x/exp/constraints"

// bstNode is the basic node in a binary search tree.
type bstNode[T constraints.Ordered] struct {
	value T

	// The two children nodes.
	left, right *bstNode[T]
}

// HasLeft reports if this node has a Left child.
func (t *bstNode[T]) HasLeft() bool {
	if t == nil {
		return false
	}

	return t.left != nil
}

// HasRight reports if this node has a Right child.
func (t *bstNode[T]) HasRight() bool {
	if t == nil {
		return false
	}

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
		// If this node is nil, the parent caller needs to create the new
		// node and set the pointer to it.
		return false
	}

	// Duplicates are not allowed.
	if v == t.value {
		return false
	}

	// If we need to go farther left, add a new node if needed,
	// otherwise recurse!
	if v < t.value {
		if t.left == nil {
			t.left = &bstNode[T]{
				value: v,
				left:  nil,
				right: nil,
			}

			return true
		}

		return t.left.Insert(v)
	}

	if t.right == nil {
		t.right = &bstNode[T]{
			value: v,
			left:  nil,
			right: nil,
		}

		return true
	}

	return t.right.Insert(v)
}

// findMin finds the node with the minimum value in the subtree rooted at t.
func (t *bstNode[T]) findMin() *bstNode[T] {
	if t == nil {
		return nil
	}
	for t.left != nil {
		t = t.left
	}

	return t
}

// deleteInternal performs the actual deletion and returns the new root.
// This is an internal method that handles root changes properly.
func (t *bstNode[T]) deleteInternal(v T) (*bstNode[T], bool) {
	if t == nil {
		return nil, false
	}

	if v < t.value {
		var deleted bool
		t.left, deleted = t.left.deleteInternal(v)

		return t, deleted
	} else if v > t.value {
		var deleted bool
		t.right, deleted = t.right.deleteInternal(v)

		return t, deleted
	}

	// Node found - handle the three deletion cases

	// Case 1: Node with no children (leaf node)
	if t.left == nil && t.right == nil {
		return nil, true
	}

	// Case 2a: Node with only right child
	if t.left == nil {
		return t.right, true
	}

	// Case 2b: Node with only left child
	if t.right == nil {
		return t.left, true
	}

	// Case 3: Node with two children
	// Replace with inorder successor (minimum value in right subtree)
	successor := t.right.findMin()
	t.value = successor.value
	// TODO(rsned): Check the success param and handle the failure case.
	t.right, _ = t.right.deleteInternal(successor.value)

	return t, true
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
// For proper deletion that can change the root, use BST.Delete() instead.
func (t *bstNode[T]) Delete(v T) bool {
	if t == nil {
		return false
	}

	// For compatibility with the BinaryTree interface, we perform deletion
	// but can't change the root. This is a limitation of the current interface.
	_, deleted := t.deleteInternal(v)

	return deleted
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

// Traverse traverses the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
//
// NOTE: Nodes in general are not expected to initiate the traverse. It would
// normally be kicked off by the main container type, e.g., BST not bstNode.
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

// Clone creates a deep copy of this BST node and its subtree.
func (t *bstNode[T]) Clone() Tree[T] {
	if t == nil {
		return nil
	}

	clone := &bstNode[T]{
		value: t.value,
		left:  nil,
		right: nil,
	}

	if t.left != nil {
		if leftClone, ok := t.left.Clone().(*bstNode[T]); ok {
			clone.left = leftClone
		}
	}

	if t.right != nil {
		if rightClone, ok := t.right.Clone().(*bstNode[T]); ok {
			clone.right = rightClone
		}
	}

	return clone
}
