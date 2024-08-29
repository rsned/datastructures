package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBSTNodeInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			tree: NewBST[int](),
			val:  5,
			want: true,
		},
		{
			tree: &BST[int]{},
			val:  5,
			want: true,
		},
		{
			// Test insert into a nil node.
			tree: (&bstNode[int]{}).left,
			val:  5,
			want: false,
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			},
			val:  5,
			want: true,
		},
		{
			// Insert to right.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			},
			val:  53,
			want: true,
		},
		{
			// Attempt to insert a duplicate value.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			},
			val:  42,
			want: false,
		},
		{
			// Insert to right -> right -> right.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 29,
							},
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
							value: 57,
						},
					},
				},
			},
			val:  85,
			want: true,
		},
		{
			// Insert to left -> right -> left -> right.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left: &bstNode[int]{
						value: 21,
						left: &bstNode[int]{
							value: 1,
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 28,
							},
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
							value: 57,
						},
					},
				},
			},
			val:  29,
			want: true,
		},
	}

	for _, test := range tests {
		if got := test.tree.Insert(test.val); got != test.want {
			t.Errorf("Insert(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestBSTNodeDelete(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			tree: NewBST[int](),
			want: false,
		},
		{
			tree: &BST[int]{},
			want: false,
		},
		{
			// Test insert into a nil node.
			tree: (&bstNode[int]{}).left,
			val:  5,
			want: false,
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		if got := test.tree.Delete(test.val); got != test.want {
			t.Errorf("Delete(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestBSTNodeSearch(t *testing.T) {
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
			tree: &BST[int]{},
			val:  5,
			want: false,
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
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
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 29,
							},
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
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

func TestBSTNodeHeight(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		want int
	}{
		{
			tree: NewBST[int](),
			want: 0,
		},
		{
			tree: &BST[int]{},
			want: 0,
		},
		{
			// Insert to left.
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
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
						},
						right: &bstNode[int]{
							value: 30,
							left: &bstNode[int]{
								value: 29,
							},
						},
					},
					right: &bstNode[int]{
						value: 84,
						left: &bstNode[int]{
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

func TestBSTNodeTraverse(t *testing.T) {
	tree := &BST[int]{
		root: &bstNode[int]{
			value: 42,
			left: &bstNode[int]{
				value: 21,
				left: &bstNode[int]{
					value: 1,
				},
				right: &bstNode[int]{
					value: 30,
					left: &bstNode[int]{
						value: 29,
					},
				},
			},
			right: &bstNode[int]{
				value: 84,
				left: &bstNode[int]{
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

func TestBSTNodeBasics(t *testing.T) {
	node := &bstNode[int]{
		value: 42,
		left: &bstNode[int]{
			value: 21,
			left: &bstNode[int]{
				value: 1,
			},
			right: &bstNode[int]{
				value: 30,
				left: &bstNode[int]{
					value: 29,
				},
			},
		},
		right: &bstNode[int]{
			value: 84,
			left: &bstNode[int]{
				value: 57,
			},
		},
	}

	if node.Right().HasRight() {
		t.Errorf("t.Right().HasRight() == true, should be false")
	}

	if node.Left().Left().HasLeft() {
		t.Errorf("t.Left().Left().HasLeft() == true, should be false")
	}

	if node.Left().Right().Left().HasLeft() {
		t.Errorf("t.Left().Left().HasLeft() == true, should be false")
	}

	r := node.Right().Right()
	if r.HasLeft() {
		t.Errorf("aaa")
	}
	if r.HasRight() {
		t.Errorf("aaa")
	}

	if node.Metadata() != "" {
		t.Errorf("There should not be any metadata on BSTs")
	}
}
