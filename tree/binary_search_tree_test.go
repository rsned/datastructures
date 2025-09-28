package tree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Most of these tests are minimal because the main testing is on the node type.

// newIntBST creates a new instance of a an integer BST.
func newIntBST() *BST[int] {
	return &BST[int]{
		root: nil,
	}
}

var (
	reasonablyFullIntBST = &BST[int]{
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
)

func TestBSTInsert(t *testing.T) {
	// Tests are done with ints to prove the code does the right thing.
	tests := []struct {
		name          string
		tree          *BST[int]
		val           int
		want          bool
		wantStructure *BST[int]
	}{
		{
			name: "Insert to empty tree",
			tree: newIntBST(),
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
			name: "Insert to left",
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
			name: "Insert to right",
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
			name: "Attempt to insert duplicate value",
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
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Insert(test.val); got != test.want {
				t.Errorf("Insert(%v) = %v, want %v", test.val, got, test.want)
			}

			if !BinaryTreesEqual(test.tree.Root(), test.wantStructure.Root()) {
				t.Errorf("value was inserted but the resulting tree was not as expected.\ngot:\n%s\nwant:\n%s\n",
					PrintBinaryTreeASCII("got", test.tree.root),
					PrintBinaryTreeASCII("want", test.wantStructure.root))
			}
		})
	}
}

func TestBSTDelete(t *testing.T) {
	tests := []struct {
		name        string
		tree        func() *BST[int]
		deleteValue int
		want        bool
		wantValues  []int // Expected in-order traversal after deletion
	}{
		{
			name:        "Delete from empty tree",
			tree:        newIntBST,
			deleteValue: 10,
			want:        false,
			wantValues:  []int{},
		},
		{
			name: "Delete non-existent value",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)

				return tree
			},
			deleteValue: 20,
			want:        false,
			wantValues:  []int{5, 10, 15},
		},
		{
			name: "Delete leaf node (left)",
			tree: func() *BST[int] {
				tree := newIntBST()

				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)

				return tree
			},
			deleteValue: 5,
			want:        true,
			wantValues:  []int{10, 15},
		},
		{
			name: "Delete leaf node (right)",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)

				return tree
			},
			deleteValue: 15,
			want:        true,
			wantValues:  []int{5, 10},
		},
		{
			name: "Delete node with only left child",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)
				tree.Insert(3)
				// Delete 5 (has only left child 3)
				return tree
			},
			deleteValue: 5,
			want:        true,
			wantValues:  []int{3, 10, 15},
		},
		{
			name: "Delete node with only right child",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)
				tree.Insert(17)
				// Delete 15 (has only right child 17)
				return tree
			},
			deleteValue: 15,
			want:        true,
			wantValues:  []int{5, 10, 17},
		},
		{
			name: "Delete node with two children",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)
				tree.Insert(3)
				tree.Insert(7)
				tree.Insert(12)
				tree.Insert(17)
				// Delete 5 (has both children 3 and 7)
				return tree
			},
			deleteValue: 5,
			want:        true,
			wantValues:  []int{3, 7, 10, 12, 15, 17},
		},
		{
			name: "Delete root node (only node)",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)

				return tree
			},
			deleteValue: 10,
			want:        true,
			wantValues:  []int{},
		},
		{
			name: "Delete root node with left child only",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)

				return tree
			},
			deleteValue: 10,
			want:        true,
			wantValues:  []int{5},
		},
		{
			name: "Delete root node with right child only",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(15)

				return tree
			},
			deleteValue: 10,
			want:        true,
			wantValues:  []int{15},
		},
		{
			name: "Delete root node with two children",
			tree: func() *BST[int] {
				tree := newIntBST()
				tree.Insert(10)
				tree.Insert(5)
				tree.Insert(15)
				tree.Insert(3)
				tree.Insert(7)
				tree.Insert(12)
				tree.Insert(17)

				return tree
			},
			deleteValue: 10,
			want:        true,
			wantValues:  []int{3, 5, 7, 12, 15, 17},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := test.tree()
			got := tree.Delete(test.deleteValue)

			if got != test.want {
				t.Errorf("Delete(%v) = %v, want %v", test.deleteValue, got, test.want)
			}

			// Verify tree structure by checking in-order traversal
			gotValues := []int{}   // Initialize as empty slice, not nil
			if tree.Height() > 0 { // Only traverse non-empty trees
				ch := tree.Traverse(TraverseInOrder)
				for value := range ch {
					gotValues = append(gotValues, value)
				}
			}

			if !cmp.Equal(gotValues, test.wantValues) {
				t.Errorf("After Delete(%v), tree traversal = %v, want %v\ndiff: %v",
					test.deleteValue, gotValues, test.wantValues, cmp.Diff(test.wantValues, gotValues))
			}

			// Verify deleted value is no longer in tree
			if test.want {
				if found := tree.Search(test.deleteValue); found {
					t.Errorf("After Delete(%v), Search(%v) = true, want false", test.deleteValue, test.deleteValue)
				}
			}
		})
	}
}

func TestBSTSearch(t *testing.T) {
	tests := []struct {
		name string
		tree Tree[int]
		val  int
		want bool
	}{
		{
			name: "Search in empty tree",
			tree: newBSTTree[int](),
			val:  5,
			want: false,
		},
		{
			name: "Search in tree with one non-matching node",
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
			name: "Search in tree with one matching node",
			tree: reasonablyFullIntBST,
			val:  57,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.tree.Search(test.val); got != test.want {
				t.Errorf("%s. Search(%v) = %v, want %v", test.name, test.val, got, test.want)
			}
		})
	}
}

func TestBSTHeight(t *testing.T) {
	tests := []struct {
		name string
		tree Tree[int]
		want int
	}{
		{
			name: "Search in empty tree",
			tree: newBSTTree[int](),
			want: 0,
		},
		{
			name: "Height of tree with one node",
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
				t.Errorf("%s. Height() = %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestBSTTraverse(t *testing.T) {
	tests := []struct {
		name  string
		tree  Tree[int]
		order TraverseOrder
		want  []int
	}{
		{
			name:  "Traverse in-order",
			tree:  reasonablyFullIntBST,
			order: TraverseInOrder,
			want:  []int{1, 21, 29, 30, 42, 57, 84},
		},
		{
			name:  "Traverse pre-order",
			tree:  reasonablyFullIntBST,
			order: TraversePreOrder,
			want:  []int{42, 21, 1, 30, 29, 84, 57},
		},
		{
			name:  "Traverse post-order",
			tree:  reasonablyFullIntBST,
			order: TraversePostOrder,
			want:  []int{1, 29, 30, 21, 57, 84, 42},
		},
		{
			name:  "Traverse reverse-order",
			tree:  reasonablyFullIntBST,
			order: TraverseReverseOrder,
			want:  []int{84, 57, 42, 30, 29, 21, 1},
		},
		{
			name:  "Traverse level-order",
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
				t.Errorf("%s. tree.Traverse() = %+v, want: %+v\ndiff: %+v",
					test.name, got, test.want, cmp.Diff(test.want, got))
			}
		})
	}
}

func TestBSTClone(t *testing.T) {
	// Create a BST tree and insert some values
	original := newIntBST()
	values := []int{5, 3, 7, 2, 4, 6, 8}

	for _, v := range values {
		original.Insert(v)
	}

	// Clone the tree
	cloned, _ := original.Clone().(*BST[int])

	// Verify that both trees are equal using BinaryTreesEqual
	if !BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Cloned BST tree should be equal to original tree")
	}

	// Test that modifications to the clone don't affect the original
	cloned.Insert(10)
	if BinaryTreesEqual(original.Root(), cloned.Root()) {
		t.Error("Original tree should not be equal to modified clone")
	}
}
