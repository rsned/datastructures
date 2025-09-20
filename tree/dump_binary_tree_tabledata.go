package tree

type indentOptionsMap map[int]indentOptions

func init() {
	var emptyOptions = indentOptions{
		nodeWidth:        0,
		prefixPadding:    0,
		intraNodePadding: 0,
		interTreePadding: 0,
		shoulderPadding:  0,
		legDepth:         0,
	}

	for nodeWidth, depths := range nodeWidthLegDepths {
		for i, d := range depths {
			prev := emptyOptions
			if o, ok := binaryTreeSpacingData[nodeWidth][i]; ok {
				prev = o
			}

			binaryTreeSpacingData[nodeWidth][i+1] = indentOptions{
				nodeWidth:        nodeWidth,
				prefixPadding:    max(0, nodeWidth+prev.legDepth+prev.shoulderPadding+prev.prefixPadding),
				intraNodePadding: nodeWidth,
				interTreePadding: max(0, prev.interTreePadding+2*(nodeWidth+prev.legDepth+prev.shoulderPadding)),
				shoulderPadding:  max(0, prev.shoulderPadding+prev.legDepth+(prev.interTreePadding-nodeWidth)/2-d),
				legDepth:         d,
			}
		}
	}
}

var (
	// nodeWidthLegDepths is a list of rendering leg depths for each level in
	// the spacing data. This will be used to generate all the indentOptions
	// instead of having to manually compute every one and redo on each
	// fine-tuning.
	nodeWidthLegDepths = map[int][]int{
		1:  {1, 1, 1, 1, 1, 1, 1, 1},
		3:  {1, 1, 1, 1, 1, 1, 1, 1},
		5:  {1, 1, 1, 1, 1, 1, 1, 1},
		7:  {1, 1, 1, 1, 1, 1, 1, 1},
		9:  {1, 1, 1, 1, 1, 1, 1, 1},
		11: {1, 1, 1, 1, 1, 1, 1, 1},
	}

	/*
		// For a more comfortable vertically spaced tree, these values are
		// visually pleasing.
			nodeWidthLegDepths = map[int][]int{
				1:  {1, 1, 2, 2, 2, 2, 2, 2},
				3:  {1, 2, 3, 3, 3, 3, 3, 3},
				5:  {1, 4, 5, 5, 5, 5, 5, 5},
				7:  {3, 4, 5, 5, 5, 5, 5, 5},
				9:  {3, 4, 5, 5, 5, 5, 5, 5},
				11: {3, 4, 5, 5, 5, 5, 5, 5},
			}
	*/

	// binaryTreeSpacingData is a mapping of node width to a map of depth from
	// bottom level being rendered to that level's indent options.
	//
	// For this data, I choose to use only odd sized width value because that
	// makes the trees look better and prevents having to deal with rounding
	// half-step spacing issues up or down.
	//
	// The data is prepopulated with the ground level spacings from which the
	// layers above are filled in during init().
	binaryTreeSpacingData = map[int]indentOptionsMap{
		1: map[int]indentOptions{
			0: {
				nodeWidth:        1,
				prefixPadding:    0,
				intraNodePadding: 1,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},

		// width 2-3
		3: map[int]indentOptions{
			0: {
				nodeWidth:        3,
				prefixPadding:    0,
				intraNodePadding: 3,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},

		// 4-5 width
		5: map[int]indentOptions{
			0: {
				nodeWidth:        5,
				prefixPadding:    0,
				intraNodePadding: 5,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},
		// 6-7 width
		7: map[int]indentOptions{
			0: {
				nodeWidth:        7,
				prefixPadding:    0,
				intraNodePadding: 7,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},
		// 8-9 width
		9: map[int]indentOptions{
			0: {
				nodeWidth:        9,
				prefixPadding:    0,
				intraNodePadding: 9,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},
		// 10-11 width
		11: map[int]indentOptions{
			0: {
				nodeWidth:        9,
				prefixPadding:    0,
				intraNodePadding: 9,
				interTreePadding: 3,
				shoulderPadding:  0,
				legDepth:         0,
			},
		},
	}
)
