package tree

import "golang.org/x/exp/constraints"

const treeTypeBST = "BST"

// BST implements a basic, unbalanced Binary Search Tree.
//
// In a BST, for each node, all values in the left subtree are less than the
// node's value, and all values in the right subtree are greater. This property
// allows for efficient searching, insertion, and deletion operations in the
// average case, but performance can degrade to O(n) in the worst case if the
// tree becomes unbalanced (e.g., when inserting sorted data).
//
// This implementation satisfies the Tree[T] interface.
type BST[T constraints.Ordered] struct {
	root *bstNode[T]
}

// NewBST creates and returns an empty Binary Search Tree, ready for use.
func NewBST[T constraints.Ordered]() Tree[T] {
	return &BST[T]{
		root: nil,
	}
}

// Root returns the root node of the tree, which implements the BinaryTree[T]
// interface. If the tree is empty, it returns nil.
func (t *BST[T]) Root() BinaryTree[T] {
	if t.root == nil {
		return nil
	}
	return t.root
}

// Insert adds a new value to the tree.
// It adheres to the BST property: values less than a node are placed in the
// left subtree, and values greater are in the right. Duplicate values are not
// inserted.
// It returns true if the insertion was successful, or false if the value
// already exists in the tree.
func (t *BST[T]) Insert(v T) bool {
	if t.root == nil {
		t.root = &bstNode[T]{
			value: v,
		}
		return true
	}
	return t.root.Insert(v)
}

// Delete removes the specified value from the tree.
// If the value is found and removed, it returns true. If the value is not in
// the tree, the tree remains unchanged and it returns false.
// The tree structure is adjusted as needed to maintain the BST property after
// deletion.
func (t *BST[T]) Delete(v T) bool {
	if t.root == nil {
		return false
	}
	var deleted bool
	// The deleteInternal method returns the potentially new root of the subtree.
	t.root, deleted = t.root.deleteInternal(v)
	return deleted
}

// Search checks if the specified value exists in the tree.
// It returns true if the value is found, and false otherwise.
func (t *BST[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}
	return t.root.Search(v)
}

// Traverse traverses the tree and returns a read-only channel of values.
// The traversal order is determined by the `tOrder` parameter.
// The channel is closed after all values have been sent.
func (t *BST[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if t.root != nil {
			traverseBinaryTree(t.root, tOrder, ch)
		}
	}()
	return ch
}

// Height returns the height of the tree, which is the number of edges on the
// longest path from the root to a leaf. An empty tree has a height of 0, and
// a tree with one node has a height of 1.
func (t *BST[T]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.Height()
}

// Clone creates and returns a deep copy of the BST. The new tree is completely
// independent of the original.
func (t *BST[T]) Clone() Tree[T] {
	if t == nil || t.root == nil {
		return NewBST[T]()
	}

	clone := &BST[T]{}
	if clonedRoot, ok := t.root.Clone().(*bstNode[T]); ok {
		clone.root = clonedRoot
	}
	return clone
}