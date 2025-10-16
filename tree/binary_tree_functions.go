package tree

import (
	"reflect"
	"slices"

	"golang.org/x/exp/constraints"
)

// BinaryTreesEquivalent checks if two binary trees contain the exact same set
// of values in the same order (i.e., their in-order traversals are identical).
// It does not compare the structural layout of the trees.
//
// For example, a skewed BST and a balanced AVL tree are equivalent if they
// contain the same elements, as their in-order traversals will produce the
// same sorted sequence.
//
// a is the first binary tree to compare.
// b is the second binary tree to compare.
// Returns true if the trees are equivalent, false otherwise.
func BinaryTreesEquivalent[T constraints.Ordered](a, b BinaryTree[T]) bool {
	if isTreeNil(a) && isTreeNil(b) {
		return true // Both are nil, so they are equivalent.
	}
	if isTreeNil(a) || isTreeNil(b) {
		return false // One is nil and the other is not.
	}

	chA := a.Traverse(TraverseInOrder)
	chB := b.Traverse(TraverseInOrder)

	for {
		valA, moreA := <-chA
		valB, moreB := <-chB

		if moreA != moreB {
			return false // One traversal finished before the other.
		}
		if !moreA {
			return true // Both traversals finished at the same time.
		}
		if valA != valB {
			return false // Values at the current position differ.
		}
	}
}

// BinaryTreesEqual checks if two binary trees are structurally identical and
// contain the same values at each corresponding node.
// This is a stricter comparison than equivalence. Both the shape of the tree
// and the values within must match.
//
// a is the first binary tree to compare.
// b is the second binary tree to compare.
// Returns true if the trees are equal, false otherwise.
func BinaryTreesEqual[T constraints.Ordered](a, b BinaryTree[T]) bool {
	// A quick check for equivalence can rule out many non-equal trees.
	// Note: This is an optimization; structural equality implies equivalence.
	if !BinaryTreesEquivalent(a, b) {
		return false
	}
	return binaryTreeStructureEqual(a, b)
}

// binaryTreeStructureEqual is an internal helper that compares the structural
// layout of two binary trees, ignoring their values.
func binaryTreeStructureEqual[T constraints.Ordered](a, b BinaryTree[T]) bool {
	aForm := getTreeStructure(a)
	bForm := getTreeStructure(b)
	return slices.Equal(aForm, bForm)
}

// getTreeStructure generates a string representation of a tree's structure
// by performing a traversal and recording the path taken (e.g., "Down-Left",
// "Visit", "Up").
func getTreeStructure[T constraints.Ordered](tree BinaryTree[T]) []string {
	var path []string
	if isTreeNil(tree) {
		return path
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		traverseBinaryTreeStructure(tree, ch)
	}()

	for p := range ch {
		path = append(path, p)
	}
	return path
}

// traverseBinaryTreeStructure is a recursive helper that walks a tree and sends
// structural path identifiers ("↓L", "↓R", "V", "↑") to a channel.
func traverseBinaryTreeStructure[T constraints.Ordered](tree BinaryTree[T], ch chan string) {
	if isTreeNil(tree) {
		return
	}

	if tree.HasLeft() {
		ch <- "↓L" // Go down to the left child
		traverseBinaryTreeStructure(tree.Left(), ch)
		ch <- "↑" // Go up from the left child
	}
	ch <- "V" // Visit node
	if tree.HasRight() {
		ch <- "↓R" // Go down to the right child
		traverseBinaryTreeStructure(tree.Right(), ch)
		ch <- "↑" // Go up from the right child
	}
}

// isTreeNil provides a safe way to check if an interface value that holds a
// tree is nil. A direct `tree == nil` check can be misleading in Go if the
// interface holds a nil pointer of a concrete type. This function uses
// reflection to inspect the underlying value, making it reliable.
//
// This is primarily intended for utility and testing functions.
func isTreeNil(a any) bool {
	if a == nil {
		return true
	}
	// Use reflection to check if the underlying value is a nil pointer.
	v := reflect.ValueOf(a)
	return v.Kind() == reflect.Ptr && v.IsNil()
}