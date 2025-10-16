package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// avlNode represents a single node within an AVL tree.
// It contains the node's value, pointers to its children and parent, and the
// balance factor (bf), which is crucial for maintaining the tree's balance.
// It is not exported and is managed by the AVL struct.
type avlNode[T constraints.Ordered] struct {
	value  T
	bf     int // The height difference between the right and left subtrees.
	parent *avlNode[T]
	left   *avlNode[T]
	right  *avlNode[T]
}

// HasLeft reports whether the node has a left child.
func (n *avlNode[T]) HasLeft() bool {
	return n.left != nil
}

// Left returns the left child of the node as a BinaryTree[T].
// It returns nil if there is no left child.
func (n *avlNode[T]) Left() BinaryTree[T] {
	return n.left
}

// HasRight reports whether the node has a right child.
func (n *avlNode[T]) HasRight() bool {
	return n.right != nil
}

// Right returns the right child of the node as a BinaryTree[T].
// It returns nil if there is no right child.
func (n *avlNode[T]) Right() BinaryTree[T] {
	return n.right
}

// Value returns the value stored at the node.
func (n *avlNode[T]) Value() T {
	return n.value
}

// Metadata returns a string representing the node's balance factor,
// e.g., "(-1)". This is used for visualization and debugging.
func (n *avlNode[T]) Metadata() string {
	return fmt.Sprintf("(%d)", n.bf)
}

// balanceFactor calculates and returns the node's balance factor, which is
// the height of the right subtree minus the height of the left subtree.
func (n *avlNode[T]) balanceFactor() int {
	if n == nil {
		return 0
	}
	return n.right.Height() - n.left.Height()
}

// insertInternal recursively inserts a new value into the subtree rooted at `n`.
// After insertion, it rebalances the tree by performing rotations if necessary.
// It returns the new root of the (potentially modified) subtree and a boolean
// indicating if the insertion occurred.
func (n *avlNode[T]) insertInternal(v T) (*avlNode[T], bool) {
	if n == nil {
		return &avlNode[T]{value: v}, true
	}

	if v == n.value {
		return n, false // Duplicate value
	}

	var inserted bool
	if v < n.value {
		n.left, inserted = n.left.insertInternal(v)
		if n.left != nil {
			n.left.parent = n
		}
	} else {
		n.right, inserted = n.right.insertInternal(v)
		if n.right != nil {
			n.right.parent = n
		}
	}

	if !inserted {
		return n, false
	}

	// Update balance factor and rebalance if necessary.
	n.bf = n.balanceFactor()

	if n.bf > 1 { // Right-heavy
		if n.right != nil && n.right.bf < 0 { // Right-Left case
			return rotateRightLeft(n), true
		}
		// Right-Right case
		return rotateLeft(n), true
	} else if n.bf < -1 { // Left-heavy
		if n.left != nil && n.left.bf > 0 { // Left-Right case
			return rotateLeftRight(n), true
		}
		// Left-Left case
		return rotateRight(n), true
	}

	return n, true
}

// Insert adds a value to the subtree rooted at the current node.
// This method is provided to satisfy the BinaryTree[T] interface.
// WARNING: This does not handle changes to the tree's root. For safe insertion,
// always call Insert from the top-level AVL struct.
func (n *avlNode[T]) Insert(v T) bool {
	if n == nil {
		return false
	}
	_, inserted := n.insertInternal(v)
	return inserted
}

// rotateLeft performs a left rotation on the subtree rooted at `node`.
// This operation is used to rebalance the tree when a node becomes right-heavy.
// It returns the new root of the rotated subtree.
func rotateLeft[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.right == nil {
		return node
	}
	newRoot := node.right
	subtree := newRoot.left

	// Perform rotation
	newRoot.left = node
	node.right = subtree

	// Update parent pointers
	newRoot.parent = node.parent
	node.parent = newRoot
	if subtree != nil {
		subtree.parent = node
	}

	// Update balance factors
	node.bf = node.balanceFactor()
	newRoot.bf = newRoot.balanceFactor()

	return newRoot
}

// rotateRight performs a right rotation on the subtree rooted at `node`.
// This operation is used to rebalance the tree when a node becomes left-heavy.
// It returns the new root of the rotated subtree.
func rotateRight[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.left == nil {
		return node
	}
	newRoot := node.left
	subtree := newRoot.right

	// Perform rotation
	newRoot.right = node
	node.left = subtree

	// Update parent pointers
	newRoot.parent = node.parent
	node.parent = newRoot
	if subtree != nil {
		subtree.parent = node
	}

	// Update balance factors
	node.bf = node.balanceFactor()
	newRoot.bf = newRoot.balanceFactor()

	return newRoot
}

// rotateRightLeft performs a double rotation to handle the "Right-Left" case.
// It consists of a right rotation on the right child, followed by a left
// rotation on the original node.
func rotateRightLeft[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.right == nil {
		return node
	}
	node.right = rotateRight(node.right)
	if node.right != nil {
		node.right.parent = node
	}
	return rotateLeft(node)
}

// rotateLeftRight performs a double rotation to handle the "Left-Right" case.
// It consists of a left rotation on the left child, followed by a right
// rotation on the original node.
func rotateLeftRight[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.left == nil {
		return node
	}
	node.left = rotateLeft(node.left)
	if node.left != nil {
		node.left.parent = node
	}
	return rotateRight(node)
}

// Delete removes a value from the subtree.
// Note: This functionality is not yet implemented for avlNode.
func (n *avlNode[T]) Delete(_ T) bool {
	return false
}

// Search recursively looks for a value in the subtree rooted at the current node.
// It returns true if the value is found, and false otherwise.
func (n *avlNode[T]) Search(v T) bool {
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
// typically initiated from the main tree structure (e.g., AVL).
func (n *avlNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
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
func (n *avlNode[T]) Height() int {
	if n == nil {
		return 0
	}
	lHeight := n.left.Height()
	rHeight := n.right.Height()
	if lHeight > rHeight {
		return lHeight + 1
	}
	return rHeight + 1
}

// toTestString is an internal helper function used for testing. It generates a
// string representation of the node structure that can be used in test cases.
func (n *avlNode[T]) toTestString(buf *bytes.Buffer, indent int) {
	const testIndents = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
	fmt.Fprintf(buf, "%svalue: %v,\n", testIndents[:indent], n.value)
	fmt.Fprintf(buf, "%sbf: %d,\n", testIndents[:indent], n.bf)

	if n.left != nil {
		buf.WriteString(testIndents[:indent] + "left: &avlNode[T]{\n")
		n.left.toTestString(buf, indent+1)
		buf.WriteString(testIndents[:indent] + "},\n")
	}

	if n.right != nil {
		buf.WriteString(testIndents[:indent] + "right: &avlNode[T]{\n")
		n.right.toTestString(buf, indent+1)
		buf.WriteString(testIndents[:indent] + "},\n")
	}
}

// Clone creates a deep copy of the node and its entire subtree.
// It returns the new node, satisfying the Tree[T] interface.
// Parent pointers are re-established during the cloning process.
func (n *avlNode[T]) Clone() Tree[T] {
	if n == nil {
		return nil
	}

	clone := &avlNode[T]{
		value: n.value,
		bf:    n.bf,
	}

	if n.left != nil {
		if leftClone, ok := n.left.Clone().(*avlNode[T]); ok {
			clone.left = leftClone
			leftClone.parent = clone
		}
	}

	if n.right != nil {
		if rightClone, ok := n.right.Clone().(*avlNode[T]); ok {
			clone.right = rightClone
			rightClone.parent = clone
		}
	}
	return clone
}