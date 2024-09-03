package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAVLNodeBasics(t *testing.T) {
	tree := &AVL[int]{}
	tree.Insert(21)
	tree.Insert(33)
	tree.Insert(1)
	tree.Insert(11)
	tree.Insert(-13)

	node := tree.root
	if !node.HasLeft() {
		t.Errorf("node.HasLeft = false, want true")
	}

	if !node.HasRight() {
		t.Errorf("node.HasRight = false, want true")
	}

	if node.Metadata() == "" {
		t.Errorf("node.Metadata should be something but was blank")
	}

	if b := node.balanceFactor(); b != -1 {
		t.Errorf("node.balanceFactor() = %d, want -1", b)
	}

	if b := node.right.right.balanceFactor(); b != 0 {
		t.Errorf("balanceFactor on nil node should be 0 but was %v", b)
	}

	// Insert into a nil node to test that path.
	node.right.right.Insert(75)
}

func TestAVLNodeInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		have    *avlNode[int]
		val     int
		want    *avlNode[int]
		success bool
	}{
		{
			// node is nil, should end up with a non-nil node.
			have: nil,
			val:  11,
			want: &avlNode[int]{
				value: 11,
				bf:    0,
			},
			success: true,
		},
		{
			// duplicate value
			have: &avlNode[int]{
				value: 5,
				bf:    0,
			},
			val: 5,
			want: &avlNode[int]{
				value: 5,
				bf:    0,
			},
			success: false,
		},
	}

	for _, test := range tests {
		if got := test.have.Insert(test.val); got != test.success {
			t.Errorf("node.Insert(%v) = %v, want %v", test.val, got, test.success)
		}

		if !binaryTreesEqual(test.have, test.want) {
			// TODO(rsned): Use dump_binary_tree here to get the two
			// trees to show.
			t.Errorf("value was inserted, but resulting tree was not correct.")
		}
	}
}

func TestAVLNodeDelete(t *testing.T) {
	tests := []struct {
		tree *avlNode[int]
		val  int
		want bool
	}{
		{
			tree: nil,
			want: false,
		},
		{
			tree: &avlNode[int]{},
			want: false,
		},
	}

	for _, test := range tests {
		if got := test.tree.Delete(test.val); got != test.want {
			t.Errorf("Delete(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestAVLNodeSearch(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		val  int
		want bool
	}{
		{
			tree: NewAVL[int](),
			val:  5,
			want: false,
		},
		{
			tree: &AVL[int]{
				root: nil,
			},
			val:  5,
			want: false,
		},
	}

	for _, test := range tests {
		if got := test.tree.Search(test.val); got != test.want {
			t.Errorf("Search(%v) = %v, want %v", test.val, got, test.want)
		}
	}
}

func TestAVLNodeHeight(t *testing.T) {
	tests := []struct {
		tree Tree[int]
		want int
	}{
		{
			tree: NewAVL[int](),
			want: 0,
		},
		{
			tree: &AVL[int]{
				root: nil,
			},
			want: 0,
		},
	}

	for _, test := range tests {
		if got := test.tree.Height(); got != test.want {
			t.Errorf("Height() = %v, want %v", got, test.want)
		}
	}
}

func TestAVLNodeTraverse(t *testing.T) {
	node := &avlNode[int]{
		value:  21,
		bf:     -1,
		parent: nil,
		left: &avlNode[int]{
			value: 1,
			bf:    0,
			left: &avlNode[int]{
				value: -13,
				bf:    0,
				left:  nil,
				right: nil,
			},
			right: &avlNode[int]{
				value: 11,
				bf:    0,
				left:  nil,
				right: nil,
			},
		},
		right: nil,
	}

	tests := []struct {
		tree  Tree[int]
		order TraverseOrder
		want  []int
	}{
		{
			tree:  node,
			order: TraverseInOrder,
			want:  []int{-13, 1, 11, 21},
		},
		{
			tree:  node,
			order: TraversePreOrder,
			want:  []int{21, 1, -13, 11},
		},
		{
			tree:  node,
			order: TraversePostOrder,
			want:  []int{-13, 11, 1, 21},
		},
		{
			tree:  node,
			order: TraverseReverseOrder,
			want:  []int{21, 11, 1, -13},
		},
		{
			tree:  node,
			order: TraverseLevelOrder,
			want:  nil,
		},
	}

	for _, test := range tests {
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
	}
}
