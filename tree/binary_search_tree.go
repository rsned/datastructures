package tree

import "golang.org/x/exp/constraints"

// BST is the simplest binary tree type. A node value and left and right
// pointers. No balancing or shuffling.
type BST[T constraints.Ordered] struct {
	root *BSTNode[T]
}

// NewBST returns an empty BST tree ready to use.
func NewBST[T constraints.Ordered]() Tree[T] {
	return &BST[T]{}
}

// Root returns the root node of the tree.
func (t *BST[T]) Root() BinaryTree[T] {
	return t.root
}

// Insert inserts the node into the tree, growing as needed, and reports
// if the operation was successful.
func (t *BST[T]) Insert(v T) bool {
	if t.root == nil {
		t.root = &BSTNode[T]{
			value: v,
		}
		return true
	}
	return t.root.Insert(v)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *BST[T]) Delete(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Delete(v)
}

// Search reports if the given value is in the tree.
func (t *BST[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Search(v)
}

// Walk traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *BST[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t.root, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *BST[T]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.Height()
}
