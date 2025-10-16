package tree

import "golang.org/x/exp/constraints"

// bstNode represents a single node within a Binary Search Tree (BST).
// It holds the node's value and pointers to its left and right children.
// It is not exported and is managed by the BST struct.
type bstNode[T constraints.Ordered] struct {
	value T
	left  *bstNode[T]
	right *bstNode[T]
}

// HasLeft reports whether the node has a left child.
func (n *bstNode[T]) HasLeft() bool {
	if n == nil {
		return false
	}
	return n.left != nil
}

// HasRight reports whether the node has a right child.
func (n *bstNode[T]) HasRight() bool {
	if n == nil {
		return false
	}
	return n.right != nil
}

// Left returns the left child of the node as a BinaryTree[T].
// It returns nil if there is no left child.
func (n *bstNode[T]) Left() BinaryTree[T] {
	if n == nil {
		return nil
	}
	return n.left
}

// Right returns the right child of the node as a BinaryTree[T].
// It returns nil if there is no right child.
func (n *bstNode[T]) Right() BinaryTree[T] {
	if n == nil {
		return nil
	}
	return n.right
}

// Value returns the value stored at the node.
func (n *bstNode[T]) Value() T {
	return n.value
}

// Metadata returns an empty string, as a standard bstNode does not have any
// specific metadata to display. This method satisfies the BinaryTree[T] interface.
func (n *bstNode[T]) Metadata() string {
	return ""
}

// Insert adds a new value into the subtree rooted at the current node.
// It maintains the BST property. Duplicates are not allowed.
// Returns true if the value was inserted, false if it already exists.
func (n *bstNode[T]) Insert(v T) bool {
	if v == n.value {
		return false // Duplicate value
	}

	if v < n.value {
		if n.left == nil {
			n.left = &bstNode[T]{value: v}
			return true
		}
		return n.left.Insert(v)
	}

	if n.right == nil {
		n.right = &bstNode[T]{value: v}
		return true
	}
	return n.right.Insert(v)
}

// findMin finds the node with the minimum value in the subtree rooted at n.
// This is a helper function used during the deletion process to find the
// in-order successor.
func (n *bstNode[T]) findMin() *bstNode[T] {
	current := n
	for current != nil && current.left != nil {
		current = current.left
	}
	return current
}

// deleteInternal is a helper that performs the actual deletion from the subtree
// rooted at `n`. It returns the new root of the (potentially modified) subtree
// and a boolean indicating if the deletion occurred.
func (n *bstNode[T]) deleteInternal(v T) (*bstNode[T], bool) {
	if n == nil {
		return nil, false
	}

	var deleted bool
	if v < n.value {
		n.left, deleted = n.left.deleteInternal(v)
		return n, deleted
	} else if v > n.value {
		n.right, deleted = n.right.deleteInternal(v)
		return n, deleted
	}

	// Value matches, this is the node to delete.

	// Case 1: Node with no children (leaf node)
	if n.left == nil && n.right == nil {
		return nil, true
	}

	// Case 2: Node with one child
	if n.left == nil {
		return n.right, true
	} else if n.right == nil {
		return n.left, true
	}

	// Case 3: Node with two children
	// Replace with in-order successor (minimum value in the right subtree).
	successor := n.right.findMin()
	n.value = successor.value
	// Delete the in-order successor from the right subtree.
	n.right, _ = n.right.deleteInternal(successor.value)
	return n, true
}

// Delete removes a value from the subtree rooted at the current node.
// This method is provided to satisfy the BinaryTree[T] interface.
// WARNING: This method cannot change the root of the tree. For safe deletion,
// always call Delete from the top-level BST struct.
func (n *bstNode[T]) Delete(v T) bool {
	if n == nil {
		return false
	}
	_, deleted := n.deleteInternal(v)
	return deleted
}

// Search recursively looks for a value in the subtree rooted at the current node.
// It returns true if the value is found, and false otherwise.
func (n *bstNode[T]) Search(v T) bool {
	if n == nil {
		return false
	}
	if v == n.value {
		return true
	}

	if v < n.value {
		return n.left.Search(v)
	}
	return n.right.Search(v)
}

// Traverse traverses the subtree rooted at this node.
// This method is provided to satisfy the BinaryTree[T] interface. Traversal is
// typically initiated from the main tree structure (e.g., BST) rather than an
// individual node.
func (n *bstNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if n != nil {
			traverseBinaryTree(n, tOrder, ch)
		}
	}()
	return ch
}

// Height calculates the height of the subtree rooted at the current node.
// A single node has a height of 1. An empty (nil) node has a height of 0.
func (n *bstNode[T]) Height() int {
	if n == nil {
		return 0
	}
	lh := n.left.Height()
	rh := n.right.Height()
	if lh > rh {
		return lh + 1
	}
	return rh + 1
}

// Clone creates a deep copy of the node and its entire subtree.
// It returns the new node, satisfying the Tree[T] interface.
func (n *bstNode[T]) Clone() Tree[T] {
	if n == nil {
		return nil
	}

	clone := &bstNode[T]{
		value: n.value,
	}

	if n.left != nil {
		if leftClone, ok := n.left.Clone().(*bstNode[T]); ok {
			clone.left = leftClone
		}
	}

	if n.right != nil {
		if rightClone, ok := n.right.Clone().(*bstNode[T]); ok {
			clone.right = rightClone
		}
	}

	return clone
}