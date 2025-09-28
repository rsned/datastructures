package tree

import (
	"testing"
)

// newIntRedBlack creates a new instance of an integer Red-Black tree.
func newIntRedBlack() *RedBlack[int] {
	return &RedBlack[int]{
		root: nil,
	}
}

func TestRedBlackClone(t *testing.T) {
	// Create a Red-Black tree and insert some values
	original := newIntRedBlack()
	values := []int{5, 3, 7, 2, 4, 6, 8}

	for _, v := range values {
		original.Insert(v)
	}

	// Clone the tree
	cloned, _ := original.Clone().(*RedBlack[int])

	// Verify that both trees are equal using BinaryTreesEqual
	if !BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Cloned Red-Black tree should be equal to original tree")
	}

	// Test that modifications to the clone don't affect the original
	cloned.Insert(10)
	if BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Original tree should not be equal to modified clone")
	}
}
