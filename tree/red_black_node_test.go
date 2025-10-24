package tree

import "testing"

func TestFindNodeRedBlack(t *testing.T) {
	tests := []struct {
		name          string
		tree          Tree[int]
		searchValue   int
		expectedFound bool
		expectedValue int
	}{
		{
			name:          "nil_root",
			tree:          emptyRedBlack,
			searchValue:   5,
			expectedFound: false,
			expectedValue: 0,
		},
		{
			name:          "single_node_found",
			tree:          singleNodeRedBlack,
			searchValue:   5,
			expectedFound: true,
			expectedValue: 5,
		},
		{
			name:          "single_node_not_found",
			tree:          singleNodeRedBlack,
			searchValue:   3,
			expectedFound: false,
			expectedValue: 0,
		},
		{
			name:          "multi_node_found_root",
			tree:          deepRedBlack,
			searchValue:   10,
			expectedFound: true,
			expectedValue: 10,
		},
		{
			name:          "multi_node_found_leaf",
			tree:          deepRedBlack,
			searchValue:   3,
			expectedFound: true,
			expectedValue: 3,
		},
		{
			name:          "multi_node_found_internal",
			tree:          deepRedBlack,
			searchValue:   15,
			expectedFound: true,
			expectedValue: 15,
		},
		{
			name:          "multi_node_not_found",
			tree:          deepRedBlack,
			searchValue:   1,
			expectedFound: false,
			expectedValue: 0,
		},
		{
			name:          "multi_node_not_found_large",
			tree:          deepRedBlack,
			searchValue:   9,
			expectedFound: false,
			expectedValue: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// verify the type assert just in case.
			rb, ok := test.tree.(*RedBlack[int])
			if !ok {
				t.Errorf("Tree is not a RedBlack tree")

				return
			}

			result := findNodeRedBlack(rb.root, test.searchValue)

			if test.expectedFound {
				if result == nil {
					t.Errorf("Expected to find node with value %v, but got nil", test.searchValue)
				} else if result.value != test.expectedValue {
					t.Errorf("Expected node value %v, but got %v", test.expectedValue, result.value)
				}
			} else if result != nil {
				t.Errorf("Expected nil result, but got node with value %v", result.value)
			}
		})
	}
}

// This test does NOT test for non-existent node since we have already found
// the node in the previous steps in the algorithm. If the node was not found
// this method isn't called.
func TestRotateLeftRedBlack(t *testing.T) {
	tests := []struct {
		name         string
		setupTree    func() *redBlackNode[int]
		rotateNode   int
		expectedRoot int
	}{
		{
			name: "nil_node_should_return_nil",
			setupTree: func() *redBlackNode[int] {
				return nil
			},
			rotateNode:   0,
			expectedRoot: 0,
		},
		{
			//    5  ->  5
			name: "node_without_right_child_should_return_same_node",
			setupTree: func() *redBlackNode[int] {
				return &redBlackNode[int]{
					value:  5,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
			},
			rotateNode:   5,
			expectedRoot: 5,
		},
		{
			// 5          10
			//  \   ->   /
			//   10     5
			name: "basic_left_rotation",
			setupTree: func() *redBlackNode[int] {
				// Create tree: 5 -> 10
				root := &redBlackNode[int]{
					value:  5,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
				right := &redBlackNode[int]{
					value:  10,
					isRed:  true,
					parent: root,
					left:   nil,
					right:  nil,
				}
				root.right = right

				return root
			},
			rotateNode:   5,
			expectedRoot: 10,
		},
		{
			// 5          10
			//  \   ->   /
			//   10     5
			//  /	     \
			// 7          7
			name: "left_rotation_with_left_subtree_becomes_left_with_right_subtree",
			setupTree: func() *redBlackNode[int] {
				// Create tree: 5 -> (10 -> 7)
				root := &redBlackNode[int]{
					value:  5,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
				right := &redBlackNode[int]{
					value:  10,
					isRed:  true,
					parent: root,
					left:   nil,
					right:  nil,
				}
				leftSubtree := &redBlackNode[int]{
					value:  7,
					isRed:  false,
					parent: right,
					left:   nil,
					right:  nil,
				}
				root.right = right
				right.left = leftSubtree

				return root
			},
			rotateNode:   5,
			expectedRoot: 10,
		},
		{
			// 3           3
			//  \     ->    \
			//   5           10
			//    \	        /
			//     10      5
			name: "left_rotation_with_parent_becomes_right_with_left_subtree",
			setupTree: func() *redBlackNode[int] {
				// Create tree: 3 -> (5 -> 10)
				parent := &redBlackNode[int]{
					value:  3,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
				root := &redBlackNode[int]{
					value:  5,
					isRed:  true,
					parent: parent,
					left:   nil,
					right:  nil,
				}
				right := &redBlackNode[int]{
					value:  10,
					isRed:  false,
					parent: root,
					left:   nil,
					right:  nil,
				}
				parent.right = root
				root.right = right

				return parent
			},
			rotateNode:   5,
			expectedRoot: 10,
		},
		{
			//   10         10
			//  /     ->    /
			// 5           7
			//  \	      /
			//   7       5
			name: "left_rotation_of_node_with_parent_and_left_child_becomes_left_leaf_with_new_parent",
			setupTree: func() *redBlackNode[int] {
				// Create tree: 10 -> (5 -> 7)
				parent := &redBlackNode[int]{
					value:  10,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
				root := &redBlackNode[int]{
					value:  5,
					isRed:  true,
					parent: parent,
					left:   nil,
					right:  nil,
				}
				right := &redBlackNode[int]{
					value:  7,
					isRed:  false,
					parent: root,
					left:   nil,
					right:  nil,
				}
				parent.left = root
				root.right = right

				return parent
			},
			rotateNode:   5,
			expectedRoot: 7,
		},
		{
			//   15         15
			//  /     ->    /
			// 10          12
			//  \         / \
			//   12      10  14
			//    \
			//     14
			name: "left_rotation_with_complex_parent_child_relationships",
			setupTree: func() *redBlackNode[int] {
				// Create tree: 15 -> (10 -> (12 -> 14))
				parent := &redBlackNode[int]{
					value:  15,
					isRed:  false,
					parent: nil,
					left:   nil,
					right:  nil,
				}
				root := &redBlackNode[int]{
					value:  10,
					isRed:  true,
					parent: parent,
					left:   nil,
					right:  nil,
				}
				right := &redBlackNode[int]{
					value:  12,
					isRed:  false,
					parent: root,
					left:   nil,
					right:  nil,
				}
				rightRight := &redBlackNode[int]{
					value:  14,
					isRed:  true,
					parent: right,
					left:   nil,
					right:  nil,
				}
				parent.left = root
				root.right = right
				right.right = rightRight

				return parent
			},
			rotateNode:   10,
			expectedRoot: 12,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := test.setupTree()

			// Find the node to rotate (if it exists)
			var nodeToRotate *redBlackNode[int]
			var originalParent *redBlackNode[int]
			if root != nil {
				nodeToRotate = findNodeRedBlack(root, test.rotateNode)
				if nodeToRotate != nil {
					originalParent = nodeToRotate.parent
				}
			}

			newRoot := rotateLeftRedBlack(nodeToRotate)

			// To keep test simpler, we use 0 to specify the node value is nil for the checks.
			if test.expectedRoot == 0 {
				if newRoot != nil {
					t.Errorf("expected nil result, but got node with value %v", newRoot.value)
				}
			} else {
				if newRoot == nil {
					t.Errorf("expected root with value %v, but got nil", test.expectedRoot)
				} else if newRoot.value != test.expectedRoot {
					t.Errorf("expected root value %v, but got %v", test.expectedRoot, newRoot.value)
				}

				// Verify parent pointers are correctly maintained
				if newRoot.parent != nil {
					// Check that parent's child pointer points to newRoot
					if newRoot.parent.left != newRoot && newRoot.parent.right != newRoot {
						t.Errorf("parent's child pointer does not point to new root")
					}

					// Additional verification: if the rotated node had a parent,
					// ensure that parent now points to the new root
					if originalParent != nil {
						if originalParent.left != newRoot && originalParent.right != newRoot {
							t.Errorf("Original parent's child pointer was not updated to point to new root. Parent: %v, Expected: %v, Got left: %v, right: %v",
								originalParent.value, newRoot.value,
								func() int {
									if originalParent.left != nil {
										return originalParent.left.value
									}

									return 0
								}(),
								func() int {
									if originalParent.right != nil {
										return originalParent.right.value
									}

									return 0
								}())
						}
					}
				}

				// Verify the rotated node's parent is correct
				if newRoot.left != nil && newRoot.left.parent != newRoot {
					t.Errorf("left child's parent pointer is incorrect")
				}
				if newRoot.right != nil && newRoot.right.parent != newRoot {
					t.Errorf("right child's parent pointer is incorrect")
				}
			}
		})
	}
}
