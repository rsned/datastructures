package tree

import "golang.org/x/exp/constraints"

// BinaryTree is the simplest tree node type.
//
// A node value and two children (left and right).
type BinaryTree[T constraints.Ordered] interface {
	Tree[T]

	// Value returns the value at this node in the tree.
	Value() T

	// HasLeft reports if this node has a Left child.
	HasLeft() bool

	// HasRight reports if this node has a Right child.
	HasRight() bool

	// Left returns the Left child, if any, of this node.
	Left() BinaryTree[T]

	// Right returns the Right child, if any, of this node.
	Right() BinaryTree[T]

	// metadata returns a metadata string, if any, for this node in the tree.
	//
	// Some examples include Balance Factor for an AVL tree, or Red/Black for a
	// node in a Red-Black tree.
	//
	// This is primarily used when rendering the tree.
	Metadata() string
}

// traverseBinaryTree is a recursive function that traverses a BinaryTree in the given
// order emitting values to the given channel.
//
// It does NOT close the channel when it is finished.
//
// Best usage is to kick this off in a goroutine.
func traverseBinaryTree[T constraints.Ordered](tree BinaryTree[T], tOrder TraverseOrder, ch chan T) {
	// What to do if the type underlying the tree is nil?
	// We can't nil check a pointer to an interface
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
		//panic("Level Order traversal not implemented")
	default:
		// TODO(rsned): There aren't other choices, so should this be
		// an error or panic as well?
	}
}
