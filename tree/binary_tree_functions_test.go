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
			a: nil,
			b: (&BST[int]{
				root: nil,
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a: (&BST[int]{
				root: nil,
			}).Root(),
			b:              nil,
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a: (&BST[int]{
				root: nil,
			}).Root(),
			b: (&BST[int]{
				root: nil,
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		// Non-empty trees.
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			b: (&BST[int]{
				root: nil,
			}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		{
			a: (&BST[int]{
				root: nil,
			}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
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
						left:  nil,
						right: nil,
					},
					right: &bstNode[int]{
						value: 53,
						left:  nil,
						right: nil,
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
						left:  nil,
						right: nil,
					},
					right: nil,
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
						left:  nil,
						right: nil,
					},
					right: &bstNode[int]{
						value: 53,
						left:  nil,
						right: nil,
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
						left:  nil,
						right: nil,
					},
					right: &bstNode[int]{
						value: 42,
						left:  nil,
						right: nil,
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
							left:  nil,
							right: nil,
						},
						right: nil,
					},
					right: nil,
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
						left:  nil,
						right: nil,
					},
					right: &bstNode[int]{
						value: 42,
						left:  nil,
						right: nil,
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
					left:  nil,
					right: nil,
				},
			}).Root(),
			b: (&AVL[int]{
				root: &avlNode[int]{
					value:  42,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: &bstNode[int]{
						value: 53,
						left:  nil,
						right: nil,
					},
				},
			}).Root(),
			b: (&AVL[int]{
				root: &avlNode[int]{
					value:  42,
					bf:     0,
					parent: nil,
					left:   nil,
					right: &avlNode[int]{
						value:  53,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
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
							left:  nil,
							right: nil,
						},
						right: nil,
					},
					right: nil,
				},
			}).Root(),
			//   21
			//  /  \
			// 1   42
			b: (&AVL[int]{
				root: &avlNode[int]{
					value:  21,
					bf:     0,
					parent: nil,
					left: &avlNode[int]{
						value:  1,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
					right: &avlNode[int]{
						value:  42,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      false,
		},
	}

	for _, test := range tests {
		if got := BinaryTreesEquivalent(test.a, test.b); got != test.wantEquivalent {
			t.Errorf("binaryTreesEquivalent(%v, %v) = %v, want %v",
				test.a, test.b, got, test.wantEquivalent)
		}
		if got := BinaryTreesEqual(test.a, test.b); got != test.wantEqual {
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
			tree: &BST[int]{
				root: nil,
			},
			vals: nil,
			want: []string{},
		},
		{
			tree: &BST[int]{
				root: nil,
			},
			vals: []int{1},
			want: []string{"V"},
		},
		{
			tree: &BST[int]{
				root: nil,
			},
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
// it's tested by TestBinaryTreeStructure.
