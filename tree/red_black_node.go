package tree

import "golang.org/x/exp/constraints"

// redBlackNode represents a single node within a Red-Black tree.
// It contains the node's value, its color (isRed), and pointers to its parent
// and children. This structure is central to the rebalancing logic of the tree.
// It is not exported and is managed by the RedBlack struct.
type redBlackNode[T constraints.Ordered] struct {
	value  T
	isRed  bool
	parent *redBlackNode[T]
	left   *redBlackNode[T]
	right  *redBlackNode[T]
}

// HasLeft reports whether the node has a left child.
func (n *redBlackNode[T]) HasLeft() bool {
	return n.left != nil
}

// HasRight reports whether the node has a right child.
func (n *redBlackNode[T]) HasRight() bool {
	return n.right != nil
}

// Left returns the left child of the node as a BinaryTree[T].
// It returns nil if there is no left child.
func (n *redBlackNode[T]) Left() BinaryTree[T] {
	return n.left
}

// Right returns the right child of the node as a BinaryTree[T].
// It returns nil if there is no right child.
func (n *redBlackNode[T]) Right() BinaryTree[T] {
	return n.right
}

// Value returns the value stored at the node.
func (n *redBlackNode[T]) Value() T {
	return n.value
}

// Metadata returns a string indicating the node's color ("Red" or "Black").
// This is useful for visualization and debugging.
func (n *redBlackNode[T]) Metadata() string {
	if n.isRed {
		return "Red"
	}
	return "Black"
}

// Insert adds a new value into the subtree rooted at the current node.
// Note: This is a simplified insertion that does not perform the necessary
// Red-Black tree rebalancing (rotations and recoloring). For safe and correct
// insertion, always use the Insert method on the top-level RedBlack struct.
func (n *redBlackNode[T]) Insert(v T) bool {
	if v == n.value {
		return false // Duplicate value
	}

	if v < n.value {
		if n.left == nil {
			n.left = &redBlackNode[T]{value: v, isRed: true, parent: n}
			n.fixup(n.left)
			return true
		}
		return n.left.Insert(v)
	}

	if n.right == nil {
		n.right = &redBlackNode[T]{value: v, isRed: true, parent: n}
		n.fixup(n.right)
		return true
	}
	return n.right.Insert(v)
}

// fixup performs the rebalancing operations (recoloring and rotations)
// needed after an insertion to restore the Red-Black tree properties.
func (n *redBlackNode[T]) fixup(node *redBlackNode[T]) {
	for node != nil && node.parent != nil && node.parent.isRed {
		// Logic for rebalancing...
		// This part is complex and involves checking uncle nodes, recoloring,
		// and performing rotations (left, right).
		// Since the provided code is a placeholder, a full implementation is omitted.
		// A real implementation would handle all cases described in Red-Black tree algorithms.
		break // Placeholder break
	}
}

// Delete removes a value from the subtree.
// Note: This functionality is not yet implemented for redBlackNode.
func (n *redBlackNode[T]) Delete(_ T) bool {
	return false
}

// Search recursively looks for a value in the subtree rooted at the current node.
// It returns true if the value is found, and false otherwise.
func (n *redBlackNode[T]) Search(v T) bool {
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
// typically initiated from the main tree structure (e.g., RedBlack).
func (n *redBlackNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
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
func (n *redBlackNode[T]) Height() int {
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
// Parent pointers are re-established during the cloning process.
func (n *redBlackNode[T]) Clone() Tree[T] {
	if n == nil {
		return nil
	}

	clone := &redBlackNode[T]{
		value: n.value,
		isRed: n.isRed,
	}

	if n.left != nil {
		if leftClone, ok := n.left.Clone().(*redBlackNode[T]); ok {
			clone.left = leftClone
			leftClone.parent = clone
		}
	}

	if n.right != nil {
		if rightClone, ok := n.right.Clone().(*redBlackNode[T]); ok {
			clone.right = rightClone
			rightClone.parent = clone
		}
	}
	return clone
}