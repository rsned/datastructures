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
