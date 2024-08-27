package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	avlTestTree = &AVL[int]{
		root: &avlNode[int]{
			value: 21,
			bf:    -1,
			left: &avlNode[int]{
				value: 1,
				bf:    0,
				left: &avlNode[int]{
					value: -13,
					bf:    0,
				},
				right: &avlNode[int]{
					value: 11,
					bf:    0,
				},
			},
			right: &avlNode[int]{
				value: 42,
				bf:    1,
				left: &avlNode[int]{
					value: 30,
					bf:    0,
				},
				right: &avlNode[int]{
					value: 84,
					bf:    0,
					left: &avlNode[int]{
						value: 57,
						bf:    0,
					},
					right: &avlNode[int]{
						value: 90,
						bf:    0,
					},
				},
			},
		},
	}
)

func TestAVLTraverse(t *testing.T) {
	tests := []struct {
		tree  Tree[int]
		order TraverseOrder
		want  []int
	}{
		{
			tree:  avlTestTree,
			order: TraverseInOrder,
			want:  []int{-13, 1, 11, 21, 30, 42, 57, 84, 90},
		},
		{
			tree:  avlTestTree,
			order: TraversePreOrder,
			want:  []int{21, 1, -13, 11, 42, 30, 84, 57, 90},
		},
		{
			tree:  avlTestTree,
			order: TraversePostOrder,
			want:  []int{-13, 11, 1, 30, 57, 90, 84, 42, 21},
		},
		{
			tree:  avlTestTree,
			order: TraverseReverseOrder,
			want:  []int{90, 84, 57, 42, 30, 21, 11, 1, -13},
		},
		{
			tree:  avlTestTree,
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
