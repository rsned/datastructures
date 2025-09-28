package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBinaryTreesEquivalentAndEqual(t *testing.T) {
	tests := []struct {
		name           string
		a, b           BinaryTree[int]
		wantEquivalent bool
		wantEqual      bool
	}{
		// Nil and empty trees.
		{
			name:           "Nil and nil",
			a:              nil,
			b:              nil,
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			name:           "Nil and empty tree",
			a:              nil,
			b:              (newIntBST()).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			name:           "Empty tree and nil",
			a:              (newIntBST()).Root(),
			b:              nil,
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			name:           "Empty tree and empty tree",
			a:              (newIntBST()).Root(),
			b:              (newIntBST()).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		// Non-empty trees.
		{
			name: "Non-empty tree and empty tree",
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			b:              (newIntBST()).Root(),
			wantEquivalent: false,
			wantEqual:      false,
		},
		{
			name: "Empty tree and non-empty tree",
			a:    (newIntBST()).Root(),
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
			name: "Two single node same type tree with same value.",
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
		{
			name: "BSTs with the different size and values.",
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
		{
			name: "BSTs with the same size and shape but different values.",
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
		{
			name: "BSTs with the same overall values but different layout.",
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
		{
			name: "AVL and BST single nodes treeswith the same value and same layout.",
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
			name: "AVL and BST multiple nodes trees with the same values and same layout.",
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
						bf:     1,
						parent: nil,
						left:   nil,
						right:  nil,
					},
				},
			}).Root(),
			wantEquivalent: true,
			wantEqual:      true,
		},
		{
			name: "AVL and BST with the same overall values but different layout.",
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
			t.Errorf("%s: binaryTreesEquivalent(%v, %v) = %v, want %v",
				test.name, test.a, test.b, got, test.wantEquivalent)
		}
		if got := BinaryTreesEqual(test.a, test.b); got != test.wantEqual {
			t.Errorf("%s: binaryTreesEqual(%v, %v) = %v, want %v",
				test.name, test.a, test.b, got, test.wantEqual)
		}
	}
}

func TestBinaryTreeStructure(t *testing.T) {
	tests := []struct {
		name string
		tree *BST[int]
		vals []int
		want []string
	}{
		{
			name: "Empty tree",
			tree: newIntBST(),
			vals: nil,
			want: []string{},
		},
		{
			name: "Tree with one node",
			tree: newIntBST(),
			vals: []int{1},
			want: []string{"V"},
		},
		{
			name: "Tree with three nodes",
			tree: newIntBST(),
			vals: []int{21, 1, 42},
			want: []string{"↓L", "V", "↑", "V", "↓R", "V", "↑"},
		},
	}

	for _, test := range tests {
		for _, val := range test.vals {
			test.tree.Insert(val)
		}

		got := binaryTreeStructure(test.tree.Root())

		if !cmp.Equal(test.want, got, cmpopts.EquateEmpty()) {
			t.Errorf("%s: binaryTreeStructure(%+v) = %+v, want %+v\ndiff: %+v",
				test.name, test.tree, got, test.want, cmp.Diff(test.want, got))
		}
	}
}

func TestBinaryTreeStructureEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b BinaryTree[int]
		want bool
	}{
		// Empty tree cases
		{
			name: "Both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "Both empty trees",
			a:    (newIntBST()).Root(),
			b:    (newIntBST()).Root(),
			want: true,
		},
		{
			name: "Nil and empty tree",
			a:    nil,
			b:    (newIntBST()).Root(),
			want: true,
		},
		{
			name: "Empty tree and nil",
			a:    (newIntBST()).Root(),
			b:    nil,
			want: true,
		},
		{
			name: "Empty BST and empty AVL",
			a:    (newIntBST()).Root(),
			b:    (newIntAVL()).Root(),
			want: true,
		},
		// One empty, one not empty
		{
			name: "Empty tree and single node tree",
			a:    (newIntBST()).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			want: false,
		},
		{
			name: "Single node tree and empty tree",
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			b:    (newIntBST()).Root(),
			want: false,
		},
		// Same type, equal structure and values
		{
			name: "Same BST type with equal structure and values",
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
			want: true,
		},
		{
			name: "Same BST type with complex equal structure and values",
			// Tree A:    21
			//          /  \
			//         1   53
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
			// Tree B:    21
			//          /  \
			//         1   53
			b: (&BST[int]{
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
			want: true,
		},
		// Different types, equal structure
		{
			name: "BST and AVL with equal structure",
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
			want: true,
		},
		{
			name: "BST and AVL with complex equal structure",
			// BST Tree:    21
			//            /  \
			//           1   53
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
			// AVL Tree:    21
			//            /  \
			//           1   53
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
						value:  53,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
				},
			}).Root(),
			want: true,
		},
		// Structure equal but values different
		{
			name: "Same structure, different values - single node",
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			}).Root(),
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 100,
					left:  nil,
					right: nil,
				},
			}).Root(),
			want: true,
		},
		{
			name: "Same structure, different values - complex tree",
			// Tree A:    21
			//          /  \
			//         1   53
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
			// Tree B:    100
			//          /    \
			//         50    200
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 100,
					left: &bstNode[int]{
						value: 50,
						left:  nil,
						right: nil,
					},
					right: &bstNode[int]{
						value: 200,
						left:  nil,
						right: nil,
					},
				},
			}).Root(),
			want: true,
		},
		{
			name: "Same structure, different values - BST vs AVL",
			// BST Tree:    21
			//            /  \
			//           1   53
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
			// AVL Tree:    100
			//            /    \
			//           50    200
			b: (&AVL[int]{
				root: &avlNode[int]{
					value:  100,
					bf:     0,
					parent: nil,
					left: &avlNode[int]{
						value:  50,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
					right: &avlNode[int]{
						value:  200,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
				},
			}).Root(),
			want: true,
		},
		// Different structures
		{
			name: "Different structures - single node vs two nodes",
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
					left: &bstNode[int]{
						value: 1,
						left:  nil,
						right: nil,
					},
					right: nil,
				},
			}).Root(),
			want: false,
		},
		{
			name: "Different structures - left-heavy vs balanced, equal values",
			// Tree A:    42
			//          /
			//         21
			//        /
			//       1
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
			// Tree B:    21
			//          /  \
			//         1   42
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
			want: false,
		},
		{
			name: "Different structures - different tree shapes",
			// Tree A:    42
			//          /  \
			//         21  53
			a: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
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
			// Tree B:    42
			//          /
			//         21
			b: (&BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left:  nil,
						right: nil,
					},
					right: nil,
				},
			}).Root(),
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := binaryTreeStructureEqual(test.a, test.b)
			if got != test.want {
				t.Errorf("binaryTreeStructureEqual(%v, %v) = %v, want %v",
					test.a, test.b, got, test.want)
			}
		})
	}
}

// traverseBinaryTreeStructure isn't tested directly since its more
// of a change detector and it's tested by TestBinaryTreeStructure.
