package tree

type indentOptionsMap map[int]indentOptions

var (
	// indentSizeLegDepths is a list of rendering leg depths for each level in the spacing data.
	// This will be used to generate all the indentOptions instead of having to manually
	// compute every one and redo on each fine-tuning.
	indentSizeLegDepths = map[int][]int{
		1: []int{0, 1, 1, 2, 2, 2, 2},
		3: []int{0, 1, 2, 3, 3, 3, 3},
		5: []int{0, 1, 4, 5, 5, 5, 5},
		7: []int{0, 3, 4, 5, 5, 5, 5},
	}

	// binaryTreeSpacingData is a mapping of node width to a map of depth from
	// bottom level being rendered to that level's indent options.
	binaryTreeSpacingData = map[int]indentOptionsMap{
		1: map[int]indentOptions{
			0: indentOptions{
				indentWidth:      1, // = nodeSize
				prefixPadding:    0, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 0 + 0 + 0 + 0 = 0
				intraNodePadding: 1, // = nodeSize
				interTreePadding: 3, // = nodeSize
				shoulderPadding:  0,
				legDepth:         0,
			},
			1: indentOptions{
				indentWidth:      1, // = nodeSize
				prefixPadding:    1, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 1 + 0 + 0 + 0 = 1
				intraNodePadding: 1, // = nodeSize
				interTreePadding: 5, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (1+0+0) = 3 + 2*1 = 3+2 = 5
				shoulderPadding:  0, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 0 + 0 + (3-1)/2 - 1 = 0+0+2/2 - 1 = 0+0+1-1 = 0
				legDepth:         1,
			},
			2: indentOptions{
				indentWidth:      1, // = nodeSize
				prefixPadding:    3, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 1 + 1 + 0 + 1 = 3
				intraNodePadding: 1, // = nodeSize
				interTreePadding: 9, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 5 + 2 * (1+1+0) = 5 + 2*2 = 5+4 = 9
				shoulderPadding:  2, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 0 + 1 + (5-1)/2 - 1 = 0+1+4/2 - 1 = 0+1+2-1 = 2
				legDepth:         1,
			},
			3: indentOptions{
				indentWidth:      1,  // = nodeSize
				prefixPadding:    7,  // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 1 + 1 + 2 + 3 = 7
				intraNodePadding: 1,  // = nodeSize
				interTreePadding: 17, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 9 + 2 * (1+1+2) = 9 + 2*4 = 9+8 = 17
				shoulderPadding:  6,  // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 2 + 1 + (9-1)/2 - 1 = 2+1+8/2 - 1 = 2+1+4-1 = 6
				legDepth:         2,
			},
			4: indentOptions{
				indentWidth:      1,  // = nodeSize
				prefixPadding:    16, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 1 + 2 + 6 + 7 = 16
				intraNodePadding: 1,  // = nodeSize
				interTreePadding: 33, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 17 + 2 * (1+2+6) = 17 + 2*9 = 15+18 = 33
				shoulderPadding:  14, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 6 + 2 + (17-1)/2 - 2 = 6+2+16/2 - 2 = 6+2+8-2 = 14
				legDepth:         2,
			},
			5: indentOptions{
				indentWidth:      1,  // = nodeSize
				prefixPadding:    33, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 1 + 2 + 14 + 16 = 33
				intraNodePadding: 1,  // = nodeSize
				interTreePadding: 47, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 31 + 2 * (1+2+14) = 31 + 2*8 = 31+16 = 47
				shoulderPadding:  29, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 14 + 2 + (31-1)/2 - 2 = 14+2+30/2 - 2 = 14+2+15-2 = 29
				legDepth:         2,
			},
		},

		// width 2-3
		3: map[int]indentOptions{
			0: indentOptions{
				indentWidth:      3, // = nodeSize
				prefixPadding:    0, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding
				intraNodePadding: 3, // = nodeSize
				interTreePadding: 3, // = nodeSize
				shoulderPadding:  0, // = prev.nodeSize + prev.shoulderPadding + prev.legDepth = 0 + 0 + 0 = 0
				legDepth:         0,
			},
			1: indentOptions{
				indentWidth:      3, // = nodeSize
				prefixPadding:    3, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 3 + 0 + 0 + 0 = 3
				intraNodePadding: 3, // = nodeSize
				interTreePadding: 9, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+0+0) = 3 + 2*3 = 3 + 6 = 9
				shoulderPadding:  0, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 0 + 0 + (3-3)/2 - 1 = 0+0+0/2 - 1 = 0+0+0-1 = 0
				legDepth:         1,
			},
			2: indentOptions{
				indentWidth:      3,  // = nodeSize
				prefixPadding:    7,  // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 3 + 1 + 0 + 3 = 7
				intraNodePadding: 3,  // = nodeSize
				interTreePadding: 17, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 9 + 2 * (3+1+0) = 9 + 2*4 = 9+8 = 17
				shoulderPadding:  2,  // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 0 + 1 + (9-3)/2 - 2 = 0+1+6/2 - 2 = 0+1+3-2 = 2
				legDepth:         2,
			},
			3: indentOptions{
				indentWidth:      3,  // = nodeSize
				prefixPadding:    14, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 3 + 2 + 2 + 7 = 14
				intraNodePadding: 3,  // = nodeSize
				interTreePadding: 31, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 17 + 2 * (3+2+2) = 3 + 2*7 = 17 + 14 = 31
				shoulderPadding:  8,  // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 2 + 2 + (17-3)/2 - 3 = 2+2+14/2 - 3 = 2+2+7-3 = 8
				legDepth:         3,
			},
			4: indentOptions{
				indentWidth:      3,  // = nodeSize
				prefixPadding:    28, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 3 + 3 + 8 + 14 = 28
				intraNodePadding: 3,  // = nodeSize
				interTreePadding: 59, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 31 + 2 * (3+3+8) = 31 + 2*14 = 31+28 = 59
				shoulderPadding:  22, // = prev.shoulder + prev.legDepth + (prev.interTree - nodeSize)/2  - legDepth = 8 + 3 + (31-3)/2 - 3 = 8+3+28/2 - 3 = 8+3+14-3 = 22
				legDepth:         3,
			},
			5: indentOptions{
				indentWidth:      3,  // = nodeSize
				prefixPadding:    56, // = nodeSize + prev.legDepth + prev.shoulderPadding + prev.prefixPadding = 3 + 3 + 22 + 28 = 56
				intraNodePadding: 3,  // = nodeSize
				interTreePadding: 87, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 59 + 2 * (3+3+22) = 31 + 2*28 = 31+56 = 87
				shoulderPadding:  50, // = prev.shoulder + prev.legDepth + (prev.interTree-nodeSize)/2 - legDepth = 22 + 3 + (59-3)/2 - 3 = 22+3+56/2 - 3 = 22+3+28-3 = 50
				legDepth:         3,
			},
		},

		// 4-5 width
		5: map[int]indentOptions{
			0: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    0,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 1, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  0,
				legDepth:         1,
			},

			1: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    0,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 1, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  0,
				legDepth:         4,
			},

			2: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    2,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 5, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  0,
				legDepth:         5,
			},
			3: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    4,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 9, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  2,
				legDepth:         5,
			},
			4: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    8,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 0, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  6,
				legDepth:         5,
			},
			5: indentOptions{
				indentWidth:      5, // = nodeSize
				prefixPadding:    8,
				intraNodePadding: 5, // = nodeSize
				interTreePadding: 1, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  7,
				legDepth:         5,
			},
		},
		// 6+ width
		7: map[int]indentOptions{
			0: indentOptions{
				indentWidth:      7, // = nodeSize
				prefixPadding:    0,
				intraNodePadding: 7, // = nodeSize
				interTreePadding: 1, // = prev.interTreePadding + 2 * (prev.nodeSize + prev.legDepth + prev.shoulderPadding) = 3 + 2 * (3+3+8) = 3 + 2*4 = 3+2 = 9
				shoulderPadding:  0,
				legDepth:         1,
			},
		},
	}
)
