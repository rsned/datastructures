package tree

import (
	"fmt"
	"testing"
)

// some sample BST tree building values.
var (
	// A fully filled 5 level BST tree.
	full5LevelBST = []int{
		42, 21, 27, 84, 11, 1, -13, 30, 57, 90, 5,
		37, 82, 88, 254, 18, 16, 20, 25, 26, 24,
		87, 100, 28, 83, 50, 89, 260, 70, 55, 47,
	}

	// Pathological left only values.
	leftLegOnly10LevelsBST = []int{
		9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	}

	// Pathological right only values.
	rightLegOnly10LevelsBST = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	// Pathological Zigzag values.  Builds a 10 level deep tree
	// going back and forth.
	zigZag10LevelsBST = []int{
		9, 1, 8, 2, 7, 3, 6, 4, 5, // 6, 7, 8, 9,
	}
)

func testBinaryInsertDump(t *testing.T) {
	t.Helper()
	tree := &BST[int]{root: nil}

	for _, val := range zigZag10LevelsBST {
		ok := tree.Insert(val)
		if !ok {
			t.Errorf("Insert(%v) = %v, want %v", val, ok, true)
		}
	}
	fmt.Printf("\n%s\n", PrintBinaryTreeASCII("", tree.Root()))

	t.Errorf("done")
}

func TestLastNonNilNode(t *testing.T) {
	// Create some test nodes for our tests
	node1 := &bstNode[int]{value: 1, left: nil, right: nil}
	node2 := &bstNode[int]{value: 2, left: nil, right: nil}
	node3 := &bstNode[int]{value: 3, left: nil, right: nil}

	tests := []struct {
		name  string
		nodes []BinaryTree[int]
		want  int
	}{
		{
			name:  "empty slice",
			nodes: []BinaryTree[int]{},
			want:  0,
		},
		{
			name:  "all nil slice",
			nodes: []BinaryTree[int]{nil, nil, nil},
			want:  0,
		},
		{
			name:  "single non-nil at index 0",
			nodes: []BinaryTree[int]{node1},
			want:  0,
		},
		{
			name:  "single non-nil at index 1",
			nodes: []BinaryTree[int]{nil, node1},
			want:  1,
		},
		{
			name:  "multiple elements with non-nil at end",
			nodes: []BinaryTree[int]{nil, nil, node1},
			want:  2,
		},
		{
			name:  "multiple elements with non-nil in middle",
			nodes: []BinaryTree[int]{nil, node1, nil, node2, nil},
			want:  3,
		},
		{
			name:  "multiple elements with non-nil at beginning",
			nodes: []BinaryTree[int]{node1, nil, nil, nil},
			want:  0,
		},
		{
			name:  "mixed nil and non-nil elements",
			nodes: []BinaryTree[int]{node1, nil, node2, nil, node3, nil, nil},
			want:  4,
		},
		{
			name:  "all non-nil elements",
			nodes: []BinaryTree[int]{node1, node2, node3},
			want:  2,
		},
		{
			name:  "non-nil at beginning and end with nils in between",
			nodes: []BinaryTree[int]{node1, nil, nil, node2},
			want:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lastNonNilNode(tt.nodes)
			if got != tt.want {
				t.Errorf("lastNonNilNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCenterString(t *testing.T) {
	tests := []struct {
		have  string
		width int
		want  string
	}{
		{
			have:  "",
			width: 0,
			want:  "",
		},
		{
			have:  "",
			width: 1,
			want:  " ",
		},
		{
			// width less than input should return input.
			have:  "a",
			width: 0,
			want:  "a",
		},
		{
			have:  "a",
			width: 1,
			want:  "a",
		},
		{
			// width only one more than input should return space after
			// input. integer division rounds down, so width 2 - len(s) = 1
			// => 1 / 2 => 0. So the front padding is [0:0] and the tail padding
			// is [0:1].
			have:  "a",
			width: 2,
			want:  "a ",
		},
		{
			have:  "a",
			width: 3,
			want:  " a ",
		},
		{
			have:  "a",
			width: 10,
			want:  "    a     ",
		},
	}

	for _, test := range tests {
		// TODO(rsned): replace spaces with the padChar
		if got := centerString(test.have, " ", test.width); got != test.want {
			t.Errorf("centerString(%q, %d) = %q, want %q", test.have, test.width, got, test.want)
		}
	}
}
