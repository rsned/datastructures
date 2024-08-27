package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBinaryTreesEquivalent(t *testing.T) {
	tests := []struct {
		a, b BinaryTree[int]
		want bool
	}{
		// Nil and empty trees.
		{
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			a:    nil,
			b:    (&BST[int]{}).Root(),
			want: true,
		},
		{
			a:    (&BST[int]{}).Root(),
			b:    nil,
			want: true,
		},
		{
			a:    (&BST[int]{}).Root(),
			b:    (&BST[int]{}).Root(),
			want: true,
		},
		// Non-empty trees.
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			b:    (&BST[int]{}).Root(),
			want: false,
		},
		{
			a: (&BST[int]{}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			want: false,
		},
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			want: true,
		},
	}

	for _, test := range tests {
		if got := binaryTreesEquivalent(test.a, test.b); got != test.want {
			t.Errorf("binaryTreesEquivalent(%v, %v) = %v, want %v",
				test.a, test.b, got, test.want)
		}
	}
}

func TestBinaryTreeStructure(t *testing.T) {

	tests := []struct {
		tree *BST[int]
		vals []int
		want []string
	}{
		{
			tree: &BST[int]{},
			vals: nil,
			want: []string{},
		},
		{
			tree: &BST[int]{},
			vals: []int{1},
			want: []string{"V"},
		},
		{
			tree: &BST[int]{},
			vals: []int{21, 1, 42},
			want: []string{"↓L", "V", "↑", "V", "↓R", "V", "↑"},
		},
	}

	for _, test := range tests {
		for _, val := range test.vals {
			test.tree.Insert(val)
		}

		// t.Errorf("Tree: \n%v", dumpBinaryTree("", test.tree.Root()))
		got := binaryTreeStructure(test.tree.Root())

		if !cmp.Equal(test.want, got, cmpopts.EquateEmpty()) {
			t.Errorf("binaryTreeStructure(%+v) = %+v, want %+v\ndiff: %+v",
				test.tree, got, test.want, cmp.Diff(test.want, got))
		}

	}
}
