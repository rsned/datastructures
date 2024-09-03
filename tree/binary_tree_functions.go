package tree

import (
	"reflect"
	"slices"

	"golang.org/x/exp/constraints"
)

// binaryTreesEquivalent tests if two BinaryTrees have the same values in the same order.
//
// As an initial pass, we start with step by step walking to see if they are the same.
func binaryTreesEquivalent[T constraints.Ordered](a, b BinaryTree[T]) bool {
	// If both are nil, then they are equivalent.
	if isTreeNil(a) == isTreeNil(b) && isTreeNil(a) {
		return true
	}

	// If one is nil and the other is not, then they are not equivalent.
	//
	// TODO(rsned):  Surprise twist! If one is nil and the other is empty,
	// should that count as equivalent?
	if isTreeNil(a) != isTreeNil(b) {
		return false
	}

	chA := a.Traverse(TraverseInOrder)
	chB := b.Traverse(TraverseInOrder)

	for {
		aVal, moreA := <-chA
		bVal, moreB := <-chB

		// Trees encountered differing values at the same step in the walk.
		if aVal != bVal {
			return false
		}

		// One tree finsished but the other has not.
		if moreA != moreB {
			return false
		}

		// Both traverses are finished and the values matched on every step.
		if !moreA && !moreB {
			return true
		}
	}
}

// binaryTreesEqual tests if two BinaryTrees have the same structure and values.
//
// TODO(rsned): Make this public method?
func binaryTreesEqual[T constraints.Ordered](a, b BinaryTree[T]) bool {
	// Test of they are equivalent first.
	return binaryTreesEquivalent(a, b) && binaryTreeStructureEqual(a, b)
}

func binaryTreeStructureEqual[T constraints.Ordered](a, b BinaryTree[T]) bool {
	aForm := binaryTreeStructure(a)
	bForm := binaryTreeStructure(b)

	return slices.Equal(aForm, bForm)
}

// binaryTreeStructure returns a string representation of the structure and
// an in order path through the given tree.
func binaryTreeStructure[T constraints.Ordered](tree BinaryTree[T]) []string {
	ch := make(chan string)
	go func() {
		traverseBinaryTreeStructure(tree, ch)
		close(ch)
	}()

	var got []string
	for {
		s, ok := <-ch
		if ok {
			got = append(got, s)
		} else {
			break
		}
	}

	return got
}

// traverseBinaryTreeStructure walks through a tree emitting directions and
// nodes to the given channel.
func traverseBinaryTreeStructure[T constraints.Ordered](tree BinaryTree[T], ch chan string) {
	if isTreeNil(tree) {
		return
	}
	if tree.HasLeft() {
		ch <- "↓L"
		traverseBinaryTreeStructure(tree.Left(), ch)
		ch <- "↑"
	}
	ch <- "V"
	if tree.HasRight() {
		ch <- "↓R"
		traverseBinaryTreeStructure(tree.Right(), ch)
		ch <- "↑"
	}

}

// isTreeNil checks if the tree generic instance the interface type is
// pointing to a nil.
//
// We can't check directly is ${interface_type} == nil, because of
// https://go.dev/doc/faq#nil_error.
//
// But we can check the Value of the type with reflection. Ideally this will
// only be used in testing and debugging code such as dumpBinaryTree and not
// in the normal code paths.
func isTreeNil(a any) bool {
	if a == nil {
		return true
	}
	// Use reflection to check if the underlying value is nil
	v := reflect.ValueOf(a)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
