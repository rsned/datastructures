package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Most of these tests are minimal because the main testing is on the node type.

func TestBSTInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			tree: &BST[int]{},
			val:  5,
			want: true,
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &BSTNode[int]{
					value: 42,
				},
			},
			val:  5,
			want: true,
		},
		{
			// Insert to right.
			tree: &BST[int]{
				root: &BSTNode[int]{
					value: 42,
				},
			},
			val:  53,
			want: true,
		},
		{
			// Attempt to insert a duplicate value.
			tree: &BST[int]{
				root: &BSTNode[int]{
					value: 42,
				},
			},
			val:  42,
			want: false,
		},
	}

	for _, test := range tests {
		if got := test.tree.Insert(test.val); got != test.want {
			t.Errorf("Insert(%v) = %v, want %v", test.val, got, test.want)
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
					root: &BSTNode[int]{
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
				root: &BSTNode[int]{
					value: 42,
				},
			},
			val:  5,
			want: false,
		},
		{
			tree: &BST[int]{
				root: &BSTNode[int]{
					value: 42,
					left: &BSTNode[int]{
						value: 21,
						left: &BSTNode[int]{
							value: 1,
						},
						right: &BSTNode[int]{
							value: 30,
							left: &BSTNode[int]{
								value: 29,
							},
						},
					},
					right: &BSTNode[int]{
						value: 84,
						left: &BSTNode[int]{
							value: 57,
						},
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
				root: &BSTNode[int]{
					value: 42,
				},
			},
			want: 1,
		},
		{
			tree: &BST[int]{
				root: &BSTNode[int]{
					value: 42,
					left: &BSTNode[int]{
						value: 21,
						left: &BSTNode[int]{
							value: 1,
						},
						right: &BSTNode[int]{
							value: 30,
							left: &BSTNode[int]{
								value: 29,
							},
						},
					},
					right: &BSTNode[int]{
						value: 84,
						left: &BSTNode[int]{
							value: 57,
						},
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
		root: &BSTNode[int]{
			value: 42,
			left: &BSTNode[int]{
				value: 21,
				left: &BSTNode[int]{
					value: 1,
				},
				right: &BSTNode[int]{
					value: 30,
					left: &BSTNode[int]{
						value: 29,
					},
				},
			},
			right: &BSTNode[int]{
				value: 84,
				left: &BSTNode[int]{
					value: 57,
				},
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
