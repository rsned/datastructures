package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBinaryTreesEquivalentAndEqual(t *testing.T) {
	tests := []struct {
		a, b           BinaryTree[int]
		wantEquivalent bool
		wantEqual      bool
	}{
		// Nil and empty trees.
		{
			a:              nil,
			b:              nil,
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a:              nil,
			b:              (&BST[int]{}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a:              (&BST[int]{}).Root(),
			b:              nil,
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a:              (&BST[int]{}).Root(),
			b:              (&BST[int]{}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		// Non-empty trees.
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			b:              (&BST[int]{}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		{
			a: (&BST[int]{}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
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
			wantEquivalent: true,
			wantEqual:      true,
		},
		// Compare BSTs with the differnet size and different values.
		{
			//   21
			//  /  \
			// 1   53
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 21,
					left: &bstNode[int]{
						value: 1,
					},
					right: &bstNode[int]{
						value: 53,
					},
				},
			}).Root(),
			//   21
			//  /
			// 1
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 21,
					left: &bstNode[int]{
						value: 1,
					},
				},
			}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		// Compare BSTs with the same size and different values.
		{
			//   21
			//  /  \
			// 1   53
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 21,
					left: &bstNode[int]{
						value: 1,
					},
					right: &bstNode[int]{
						value: 53,
					},
				},
			}).Root(),
			//   21
			//  /  \
			// 1   42
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 21,
					left: &bstNode[int]{
						value: 1,
					},
					right: &bstNode[int]{
						value: 42,
					},
				},
			}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		// Compare BSTs with the same overall values but different layout.
		{
			//     42
			//    /
			//   21
			//  /
			// 1
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
						},
					},
				},
			}).Root(),
			//   21
			//  /  \
			// 1   42
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 21,
					left: &bstNode[int]{
						value: 1,
					},
					right: &bstNode[int]{
						value: 42,
					},
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      false,
		},
		// Compare an AVL and BST with the same overall values and same layout.
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			}).Root(),
			b: (&AVL[int]{
				root: &avlNode[int]{
					value: 42,
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					right: &bstNode[int]{
						value: 53,
					},
				},
			}).Root(),
			b: (&AVL[int]{
				root: &avlNode[int]{
					value: 42,
					right: &avlNode[int]{
						value: 53,
					},
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		// Compare an AVL and BST with the same overall values but different layout.
		{
			//     42
			//    /
			//   21
			//  /
			// 1
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
						},
					},
				},
			}).Root(),
			//   21
			//  /  \
			// 1   42
			b: (&AVL[int]{
				root: &avlNode[int]{
					value: 21,
					left: &avlNode[int]{
						value: 1,
					},
					right: &avlNode[int]{
						value: 42,
					},
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      false,
		},
	}

	for _, test := range tests {
		if got := binaryTreesEquivalent(test.a, test.b); got != test.wantEquivalent {
			t.Errorf("binaryTreesEquivalent(%v, %v) = %v, want %v",
				test.a, test.b, got, test.wantEquivalent)
		}
		if got := binaryTreesEqual(test.a, test.b); got != test.wantEqual {
			t.Errorf("binaryTreesEqual(%v, %v) = %v, want %v",
				test.a, test.b, got, test.wantEqual)
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

// traverseBinaryTreeStructure isnt tested directly since its more of a change detector and
// and it's tested by TestBinaryTreeStructure.
