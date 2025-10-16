package tree

import (
	"golang.org/x/exp/constraints"
)

// TraverseOrder represents the common orders in which tree nodes may be
// traversed.
type TraverseOrder int

const (
	// TraverseInOrder traverses a binary tree from the left subtree to the root,
	// then to the right subtree. For a Binary Search Tree, this visits the nodes
	// in ascending sorted order.
	TraverseInOrder TraverseOrder = iota
	// TraversePreOrder traverses a binary tree from the root to the left subtree,
	// then to the right subtree.
	TraversePreOrder
	// TraversePostOrder traverses a binary tree from the left subtree to the
	// right subtree, then to the root.
	TraversePostOrder

	// TraverseReverseOrder traverses a binary tree from the right subtree to the
	// root, then to the left subtree. For a Binary Search Tree, this visits
	// the nodes in descending sorted order.
	TraverseReverseOrder

	// TraverseLevelOrder performs a breadth-first traversal where nodes are
	// visited level by level, from left to right within each level.
	TraverseLevelOrder
)

// String returns a human-readable string label for the TraverseOrder type.
func (t TraverseOrder) String() string {
	switch t {
	case TraverseInOrder:
		return "In-Order"
	case TraversePreOrder:
		return "Pre-Order"
	case TraversePostOrder:
		return "Post-Order"
	case TraverseReverseOrder:
		return "Reverse-Order"
	case TraverseLevelOrder:
		return "Level-Order"
	default:
		return "invalid traverse order"
	}
}

// Traverser is an interface for tree types that support traversal.
// It provides a standard way to iterate over the values stored in a tree.
type Traverser[T constraints.Ordered] interface {
	// Traverse traverses the tree in the specified order and returns a
	// read-only channel of values. The channel is closed after all values
	// have been sent. This allows for easy iteration over the tree's elements:
	//
	//   for value := range myTree.Traverse(tree.TraverseInOrder) {
	//       fmt.Println(value)
	//   }
	//
	// The traversal is performed concurrently, so the caller does not need to
	// wait for the entire traversal to complete before processing values.
	Traverse(order TraverseOrder) <-chan T
}

// Tree defines the basic interface common to all tree data structures in this
// package. It ensures a consistent API for fundamental tree operations.
type Tree[T constraints.Ordered] interface {
	// Insert adds a new value to the tree.
	// It returns true if the value was successfully inserted, and false if the
	// value already exists in the tree or could not be added.
	Insert(v T) bool

	// Delete removes the specified value from the tree.
	// It returns true if the value was found and removed, and false otherwise.
	// If the deleted node is not a leaf, the tree's internal structure may be
	// reorganized to maintain its properties.
	Delete(v T) bool

	// Search checks if a given value exists in the tree.
	// It returns true if the value is found, and false otherwise.
	Search(v T) bool

	// Height returns the height of the tree, which is the number of edges on
	// the longest path from the root node to a leaf node. A tree with only a
	// root node has a height of 0.
	Height() int

	// Clone creates and returns a deep copy of the entire tree. The new tree is
	// completely independent of the original.
	Clone() Tree[T]

	// Traverser is embedded, providing the Traverse method.
	Traverser[T]
}