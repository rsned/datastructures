package tree

import "golang.org/x/exp/constraints"

// BinaryTree defines the interface for a binary tree node. It extends the base
// Tree interface with methods specific to nodes in a binary tree structure,
// such as accessing its value and children.
type BinaryTree[T constraints.Ordered] interface {
	// Tree is the embedded base interface, providing common tree operations.
	Tree[T]

	// Value returns the value stored at the current node in the tree.
	Value() T

	// HasLeft reports whether the current node has a left child.
	HasLeft() bool

	// HasRight reports whether the current node has a right child.
	HasRight() bool

	// Left returns the left child of the current node. It returns nil if the
	// node does not have a left child.
	Left() BinaryTree[T]

	// Right returns the right child of the current node. It returns nil if the
	// node does not have a right child.
	Right() BinaryTree[T]

	// Metadata returns a string containing metadata about the node, which is
	// useful for visualization and debugging. For example, an AVL tree node
	// might return its balance factor, while a Red-Black tree node might
	// return its color.
	Metadata() string
}

// traverseBinaryTree is an internal recursive helper function that traverses a
// BinaryTree in a specified order and sends the node values to a channel.
//
// This function does NOT close the channel when it completes. The caller is
// responsible for managing the channel's lifecycle. It is intended to be run
// in a separate goroutine to enable concurrent traversal.
//
// tree is the starting node for the traversal.
// tOrder specifies the traversal order (e.g., In-Order, Pre-Order).
// ch is the channel to which the traversed values are sent.
func traverseBinaryTree[T constraints.Ordered](tree BinaryTree[T], tOrder TraverseOrder, ch chan T) {
	// A nil check on an interface in Go is tricky. If a concrete type with a
	// nil value is passed, the interface itself is not nil. We rely on the
	// calling context to ensure a valid tree is passed.
	switch tOrder {
	case TraverseInOrder:
		if tree.HasLeft() {
			traverseBinaryTree(tree.Left(), tOrder, ch)
		}
		ch <- tree.Value()
		if tree.HasRight() {
			traverseBinaryTree(tree.Right(), tOrder, ch)
		}
	case TraversePreOrder:
		ch <- tree.Value()
		if tree.HasLeft() {
			traverseBinaryTree(tree.Left(), tOrder, ch)
		}
		if tree.HasRight() {
			traverseBinaryTree(tree.Right(), tOrder, ch)
		}
	case TraversePostOrder:
		if tree.HasLeft() {
			traverseBinaryTree(tree.Left(), tOrder, ch)
		}
		if tree.HasRight() {
			traverseBinaryTree(tree.Right(), tOrder, ch)
		}
		ch <- tree.Value()
	case TraverseReverseOrder:
		if tree.HasRight() {
			traverseBinaryTree(tree.Right(), tOrder, ch)
		}
		ch <- tree.Value()
		if tree.HasLeft() {
			traverseBinaryTree(tree.Left(), tOrder, ch)
		}
	case TraverseLevelOrder:
	// Level-order traversal is typically handled iteratively with a queue,
	// not recursively, so it's omitted here. The public Traverse method on
	// tree implementations will handle this case.
	default:
		// Invalid traversal order; do nothing.
	}
}