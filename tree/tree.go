package tree

import (
	"golang.org/x/exp/constraints"
)

// TraverseOrder represents the common orders that tree nodes may be traversed in.
type TraverseOrder int

const (
	// TraverseInOrder traverses from the left subtree to the root then to the right subtree.
	TraverseInOrder TraverseOrder = iota
	// TraversePreOrder traverses from the root to the left subtree then to the right subtree.
	TraversePreOrder
	// TraversePostOrder traverses from the left subtree to the right subtree then to the root.
	TraversePostOrder

	// TraverseReverseOrder traverses from the right subtree to the root then to the left.
	TraverseReverseOrder

	// TraverseLevelOrder performs breadth first search where nodes are visited by level
	// left to right before moving to the next level down.
	TraverseLevelOrder

	// TODO(rsned): Are there any other reasonable paths to take that should be added?
)

// String returns a string label for the TraverseOrder type.
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

// Traverser is an interface for trees that implement a way to traverse themselves.
type Traverser[T constraints.Ordered] interface {
	// Traverse traverse the tree in the specified order emitting the values to
	// the channel. Channel is closed once the final value is emitted.
	Traverse(TraverseOrder) <-chan T
}

// Tree defines the basic interface common to all trees.
type Tree[T constraints.Ordered] interface {
	// Insert adds the given value into the true.
	// If the value could not be added, false is returned.
	Insert(v T) bool

	// Delete the requested node from the tree and reports if it was successful.
	// If the value is not in the tree, the tree is unchanged and false is returned.
	//
	// If the node is not a leaf the trees internal structure may be updated.
	Delete(v T) bool

	// Search reports if the given value is in the tree.
	Search(v T) bool

	// Height returns the height of the longest path in the tree from the
	// root node to the farthest leaf.
	Height() int

	Traverser[T]
}
