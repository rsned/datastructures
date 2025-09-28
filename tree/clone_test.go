package tree

import (
	"testing"
)

// Basic clone tests for general tree functionality.
// Tree-specific clone tests are in their respective test files (avl_test.go,
// binary_search_tree_test.go, red_black_test.go, etc.).

func TestBasicClone(t *testing.T) {
	// Test the base Clone function using just the plain BST.

	// Create a BST tree and insert some values
	original := newIntBST()
	values := []int{5, 3, 7, 2, 4, 6, 8}

	for _, v := range values {
		original.Insert(v)
	}

	// Clone the tree
	cloned, _ := original.Clone().(*BST[int])

	// Verify that both trees are equal using BinaryTreesEqual
	if !BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Cloned BST tree should be equal to original tree")
	}

	// Test that modifications to the clone don't affect the original
	cloned.Insert(10)
	if BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Original tree should not be equal to modified clone")
	}
}

func TestEmptyTreeClone(t *testing.T) {
	// Test cloning empty trees
	emptyAVL, _ := NewAVL[int]().(*AVL[int])
	clonedAVL, _ := emptyAVL.Clone().(*AVL[int])
	if !BinaryTreesEqual(emptyAVL.Root(), clonedAVL.Root()) {
		t.Error("Cloned empty AVL tree should be equal to original")
	}

	emptyBST, _ := NewBST[int]().(*BST[int])
	clonedBST, _ := emptyBST.Clone().(*BST[int])
	if !BinaryTreesEqual(emptyBST.Root(), clonedBST.Root()) {
		t.Error("Cloned empty BST tree should be equal to original")
	}

	emptyRB, _ := NewRedBlack[int]().(*RedBlack[int])
	clonedRB, _ := emptyRB.Clone().(*RedBlack[int])
	if !BinaryTreesEqual(emptyRB.Root(), clonedRB.Root()) {
		t.Error("Cloned empty Red-Black tree should be equal to original")
	}
}
