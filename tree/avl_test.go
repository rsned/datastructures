package tree

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	avlTestTree = &AVL[int]{
		root: &avlNode[int]{
			value:  21,
			bf:     -1,
			parent: nil,
			left: &avlNode[int]{
				value:  1,
				bf:     0,
				parent: nil,
				left: &avlNode[int]{
					value:  -13,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
				right: &avlNode[int]{
					value:  11,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
			},
			right: &avlNode[int]{
				value:  42,
				bf:     1,
				parent: nil,
				left: &avlNode[int]{
					value:  30,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
				right: &avlNode[int]{
					value:  84,
					bf:     0,
					parent: nil,
					left: &avlNode[int]{
						value:  57,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
					},
					right: &avlNode[int]{
						value:  90,
						bf:     0,
						parent: nil,
						left:   nil,
						right:  nil,
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

func TestAVLInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		have    *AVL[int]
		val     int
		want    *AVL[int]
		success bool
	}{
		{
			have: &AVL[int]{
				root: nil,
			},
			val: 11,
			want: &AVL[int]{
				root: &avlNode[int]{
					value:  11,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
			},
			success: true,
		},
		{
			// duplicate value
			have: &AVL[int]{
				root: &avlNode[int]{
					value:  5,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
			},
			val: 5,
			want: &AVL[int]{
				root: &avlNode[int]{
					value:  5,
					bf:     0,
					parent: nil,
					left:   nil,
					right:  nil,
				},
			},
			success: false,
		},
		{
			// right heavy node that inserting should force a re-balance.
			have: func() *AVL[int] {
				root := &avlNode[int]{
					value:  5,
					bf:     1,
					parent: nil,
					left:   nil,
					right: &avlNode[int]{
						value:  11,
						bf:     0,
						left:   nil,
						right:  nil,
						parent: nil, // will be set below
					},
				}
				root.right.parent = root

				return &AVL[int]{root: root}
			}(),
			val: 13,
			want: func() *AVL[int] {
				root := &avlNode[int]{
					value:  11,
					bf:     0,
					parent: nil,
					left: &avlNode[int]{
						value:  5,
						bf:     0,
						left:   nil,
						right:  nil,
						parent: nil, // will be set below
					},
					right: &avlNode[int]{
						value:  13,
						bf:     0,
						left:   nil,
						right:  nil,
						parent: nil, // will be set below
					},
				}
				root.left.parent = root
				root.right.parent = root

				return &AVL[int]{root: root}
			}(),
			success: true,
		},
	}

	for _, test := range tests {
		if got := test.have.Insert(test.val); got != test.success {
			t.Errorf("node.Insert(%v) = %v, want %v", test.val, got, test.success)
		}

		if !BinaryTreesEqual(test.have.Root(), test.want.Root()) {
			// TODO(rsned): Use dump_binary_tree here to get the two
			// trees to show.
			t.Errorf("value was inserted, but resulting tree was not correct.\ngot:\n%s\nwant:\n%s\n",
				PrintBinaryTreeASCII("got", test.have.Root()),
				PrintBinaryTreeASCII("want", test.want.Root()))
		}
	}
}

func testAVLInsertDump(t *testing.T) {
	t.Helper()
	tree := &AVL[int]{
		root: nil,
	}

	for _, val := range []int{
		5, 11, 13,
		6, 5, 4,
		//  3, 2, 1, 7, 8, 9,
		/*
			300, 600, 200, 700, 400,
			550, 100, 800, 250, 650,
				575, 900, 50, 150, 350,
				450, 775, 25, 75, 125,
				175, 225, 275, 350, 450,
				213, 237, 264, 283, 325,
				375, 625, 675, 10,
		*/
	} {
		tree.Insert(val)
		fmt.Printf("\n%s\n", PrintBinaryTreeASCII("Initial", tree.Root()))
	}

	t.Errorf("done: %v\n", tree.toTestString())
}
