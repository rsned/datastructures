package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Most of these tests are minimal because the main testing is on the node type.

func TestBSTInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		tree          *BST[int]
		val           int
		want          bool
		wantStructure *BST[int]
	}{
		{
			tree: &BST[int]{root: nil},
			val:  5,
			want: true,
			wantStructure: &BST[int]{
				root: &bstNode[int]{
					value: 5,
					left:  nil,
					right: nil,
				},
			},
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			val:  5,
			want: true,
			wantStructure: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 5,
						left:  nil,
						right: nil,
					},
					right: nil,
				},
			},
		},
		{
			// Insert to right.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			val:  53,
			want: true,
			wantStructure: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: &bstNode[int]{
						value: 53,
						left:  nil,
						right: nil,
					},
				},
			},
		},
		{
			// Attempt to insert a duplicate value.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			val:  42,
			want: false,
			wantStructure: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
		},
	}

	for _, test := range tests {
		if got := test.tree.Insert(test.val); got != test.want {
			t.Errorf("Insert(%v) = %v, want %v", test.val, got, test.want)
		}

		if !BinaryTreesEqual(test.tree.Root(), test.wantStructure.Root()) {
			t.Errorf("value was inserted but the resulting tree was not as expected.\ngot:\n%s\nwant:\n%s\n",
				PrintBinaryTreeASCII("got", test.tree.root),
				PrintBinaryTreeASCII("want", test.wantStructure.root))
		}
	}
}

func TestBSTDelete(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			// tree has no root node to start with.
			tree: NewBST[int](),
			val:  0,
			want: false,
		},
		{
			// Value not in tree.
			tree: NewBST[int](),
			val:  5,
			want: false,
		},
		/*
			TODO(rsned): Once Delete is implemented, add this case.
			{
				tree: &BST[int]{
					root: &bstNode[int]{
						value: 42,
					},
				}
				val: 42,
				want: false,
			},
		*/
	}

	for _, test := range tests {
		if got := test.tree.Delete(test.val); got != test.want {
			t.Errorf("Delete(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestBSTSearch(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			tree: NewBST[int](),
			val:  5,
			want: false,
		},
		{
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			val:  5,
			want: false,
		},
		{
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
							left:  nil,
							right: nil,
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 29,
								left:  nil,
								right: nil,
							},
							right: nil,
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
							value: 57,
							left:  nil,
							right: nil,
						},
						right: nil,
					},
				},
			},
			val:  57,
			want: true,
		},
	}

	for _, test := range tests {
		if got := test.tree.Search(test.val); got != test.want {
			t.Errorf("Search(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestBSTHeight(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		want int
	}{
		{
			tree: NewBST[int](),
			want: 0,
		},
		{
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			want: 1,
		},
		{
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
							left:  nil,
							right: nil,
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 29,
								left:  nil,
								right: nil,
							},
							right: nil,
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
							value: 57,
							left:  nil,
							right: nil,
						},
						right: nil,
					},
				},
			},
			want: 4,
		},
	}

	for _, test := range tests {
		if got := test.tree.Height(); got != test.want {
			t.Errorf("Height() = %v, want %v", got, test.want)
		}
	}
}

func TestBSTTraverse(t *testing.T) {
	tree := &BST[int]{
		root: &bstNode[int]{
			value: 42,
			left: &bstNode[int]{
				value: 21,
				left: &bstNode[int]{
					value: 1,
					left:  nil,
					right: nil,
				},
				right: &bstNode[int]{
					value: 30,
					left: &bstNode[int]{
						value: 29,
						left:  nil,
						right: nil,
					},
					right: nil,
				},
			},
			right: &bstNode[int]{
				value: 84,
				left: &bstNode[int]{
					value: 57,
					left:  nil,
					right: nil,
				},
				right: nil,
			},
		},
	}

	tests := []struct {
		tree  Tree[int]
		order TraverseOrder
		want  []int
	}{
		{
			tree:  tree,
			order: TraverseInOrder,
			want:  []int{1, 21, 29, 30, 42, 57, 84},
		},
		{
			tree:  tree,
			order: TraversePreOrder,
			want:  []int{42, 21, 1, 30, 29, 84, 57},
		},
		{
			tree:  tree,
			order: TraversePostOrder,
			want:  []int{1, 29, 30, 21, 57, 84, 42},
		},
		{
			tree:  tree,
			order: TraverseReverseOrder,
			want:  []int{84, 57, 42, 30, 29, 21, 1},
		},
		{
			tree:  tree,
			order: TraverseLevelOrder,
			want:  nil,
		},
	}

	for _, test := range tests {
		ch := tree.Traverse(test.order)

		var got []int
		for {
			j, ok := <-ch
			if ok {
				got = append(got, j)
			} else {
				break
			}
		}

		if !cmp.Equal(got, test.want) {
			t.Errorf("tree.Traverse() = %+v, want: %+v\ndiff: %+v",
				got, test.want, cmp.Diff(test.want, got))
		}
	}
}
