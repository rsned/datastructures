package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBSTNodeInsert(t *testing.T) {
	x, _ := reasonablyFullIntBST.Clone().(*BST[int])
	bushy := x.root

	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		name string
		tree BinaryTree[int]
		val  int
		want bool
	}{
		// We don't test the case for when the root node is nil here because
		// that is handled by the main BST Insert method (as it requires
		// replacing the root node pointer).
		{
			name: "Insert to left of root",
			tree: &bstNode[int]{
				left:  nil,
				right: nil,
				value: 42,
			},
			val:  5,
			want: true,
		},
		{
			name: "Insert to right of root",
			tree: &bstNode[int]{
				left:  nil,
				right: nil,
				value: 42,
			},
			val:  53,
			want: true,
		},
		{
			name: "Attempt to insert duplicate value",
			tree: &bstNode[int]{
				left:  nil,
				right: nil,
				value: 42,
			},
			val:  42,
			want: false,
		},
		{
			name: "Insert to bushy tree",
			tree: bushy,
			val:  85,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Insert(test.val); got != test.want {
				t.Errorf("Insert(%v) = %v, want %v", test.val, got, test.want)
			}
		})
	}
}

func TestBSTNodeDelete(t *testing.T) {
	tests := []struct {
		name string
		tree Tree[int]
		val  int
		want bool
	}{
		{
			name: "Delete from empty tree",
			tree: &BST[int]{root: nil},
			val:  0,
			want: false,
		},
		{
			name: "Delete non-existent value from tree with no child nodes",
			tree: (&bstNode[int]{value: 0, left: nil, right: nil}).left,
			val:  5,
			want: false,
		},
		{
			name: "Delete non-existent value from tree with one child nodes",
			tree: &BST[int]{
				root: &bstNode[int]{
					value: 42,
					left:  nil,
					right: nil,
				},
			},
			val:  0,
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Delete(test.val); got != test.want {
				t.Errorf("Delete(%v) = %v, want %v", test.val, got, test.want)
			}
		})
	}
}

func TestBSTNodeSearch(t *testing.T) {
	tests := []struct {
		name string
		tree Tree[int]
		val  int
		want bool
	}{
		{
			name: "Search in empty tree",
			tree: NewBST[int](),
			val:  5,
			want: false,
		},
		{
			name: "Search in tree with no child nodes",
			tree: &BST[int]{root: nil},
			val:  5,
			want: false,
		},
		{
			name: "Search non-existent value in tree with one child nodes",
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
			name: "Search existing value in tree with child nodes",
			tree: reasonablyFullIntBST.root,
			val:  57,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Search(test.val); got != test.want {
				t.Errorf("Search(%v) = %v, want %v", test.val, got, test.want)
			}
		})
	}
}

func TestBSTNodeHeight(t *testing.T) {
	tests := []struct {
		name string
		tree Tree[int]
		want int
	}{
		{
			name: "Height of empty tree",
			tree: NewBST[int](),
			want: 0,
		},
		{
			name: "Height of tree with no child nodes",
			tree: &BST[int]{root: nil},
			want: 0,
		},
		{
			name: "Height of tree with one child nodes",
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
			name: "Height of reasonably full tree",
			tree: reasonablyFullIntBST,
			want: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Height(); got != test.want {
				t.Errorf("Height() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestBSTNodeTraverse(t *testing.T) {
	tests := []struct {
		name  string
		tree  Tree[int]
		order TraverseOrder
		want  []int
	}{
		{
			name:  "Traverse in order",
			tree:  reasonablyFullIntBST,
			order: TraverseInOrder,
			want:  []int{1, 21, 29, 30, 42, 57, 84},
		},
		{
			name:  "Traverse pre order",
			tree:  reasonablyFullIntBST,
			order: TraversePreOrder,
			want:  []int{42, 21, 1, 30, 29, 84, 57},
		},
		{
			name:  "Traverse post order",
			tree:  reasonablyFullIntBST,
			order: TraversePostOrder,
			want:  []int{1, 29, 30, 21, 57, 84, 42},
		},
		{
			name:  "Traverse reverse order",
			tree:  reasonablyFullIntBST,
			order: TraverseReverseOrder,
			want:  []int{84, 57, 42, 30, 29, 21, 1},
		},
		{
			name:  "Traverse level order",
			tree:  reasonablyFullIntBST,
			order: TraverseLevelOrder,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := test.tree.Traverse(test.order)

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
		})
	}
}

func TestBSTNodeBasics(t *testing.T) {
	node := reasonablyFullIntBST.root

	if node == nil {
		t.Errorf("node should not be nil")
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
		t.Errorf("Node with no children should not have a left child")
	}
	if r.HasRight() {
		t.Errorf("Node with no children should not have a right child")
	}

	if node.Metadata() != "" {
		t.Errorf("There should not be any metadata on BSTs")
	}
}
