package tree

import (
	"math"
	"testing"
)

// Common test tree variables for various tests

// Empty trees - one of each type
var (
	emptyBST      = NewBST[int]()
	emptyAVL      = NewAVL[int]()
	emptyRedBlack = NewRedBlack[int]()
)

// Simple trees - one of each type with 3 nodes
var (
	simpleBST      Tree[int]
	simpleAVL      Tree[int]
	simpleRedBlack Tree[int]
)

// Deep trees - one of each type with at least 3 levels (height >= 2)
var (
	deepBST      Tree[int]
	deepAVL      Tree[int]
	deepRedBlack Tree[int]
)

// TODO(rsned): Would it be beneficial to generate a much bushier, deeper set
// of trees to add even more comprehensive and exhaustive tests and benchmarks.

// init creates the various test trees
func init() {
	// Simple Tree: 3 nodes
	//       10
	//      /  \
	//     5    15

	simpleBST = NewBST[int]()
	simpleBST.Insert(10)
	simpleBST.Insert(5)
	simpleBST.Insert(15)

	simpleAVL = NewAVL[int]()
	simpleAVL.Insert(5)
	simpleAVL.Insert(10)
	simpleAVL.Insert(15)

	simpleRedBlack = NewRedBlack[int]()
	simpleRedBlack.Insert(5)
	simpleRedBlack.Insert(10)
	simpleRedBlack.Insert(15)

	// Deep Tree: 6 nodes with 3 levels
	//       10
	//      /  \
	//     5    15
	//    / \     \
	//   3   7     20

	deepBST = NewBST[int]()
	deepBST.Insert(10)
	deepBST.Insert(5)
	deepBST.Insert(15)
	deepBST.Insert(3)
	deepBST.Insert(7)
	deepBST.Insert(20)

	deepAVL = NewAVL[int]()
	deepAVL.Insert(3)
	deepAVL.Insert(5)
	deepAVL.Insert(7)
	deepAVL.Insert(10)
	deepAVL.Insert(15)
	deepAVL.Insert(20)

	deepRedBlack = NewRedBlack[int]()
	deepRedBlack.Insert(3)
	deepRedBlack.Insert(5)
	deepRedBlack.Insert(7)
	deepRedBlack.Insert(10)
	deepRedBlack.Insert(15)
	deepRedBlack.Insert(20)
}

func TestCalculateBasicMetrics(t *testing.T) {
	tests := []struct {
		name       string
		tree       Tree[int]
		treeType   string
		wantMin    int
		wantMax    int
		wantHeight int
		wantNodes  int
		wantValues bool
	}{
		{
			name:       "EmptyBST",
			tree:       emptyBST,
			treeType:   treeTypeBST,
			wantMin:    0,
			wantMax:    0,
			wantHeight: 0,
			wantNodes:  0,
			wantValues: false,
		},
		{
			name:       "EmptyAVL",
			tree:       emptyAVL,
			treeType:   treeTypeAVL,
			wantMin:    0,
			wantMax:    0,
			wantHeight: 0,
			wantNodes:  0,
			wantValues: false,
		},
		{
			name:       "EmptyRedBlack",
			tree:       emptyRedBlack,
			treeType:   treeTypeRedBlack,
			wantMin:    0,
			wantMax:    0,
			wantHeight: 0,
			wantNodes:  0,
			wantValues: false,
		},
		{
			name:       "SimpleBST",
			tree:       simpleBST,
			treeType:   treeTypeBST,
			wantMin:    5,
			wantMax:    15,
			wantHeight: 2,
			wantNodes:  3,
			wantValues: true,
		},
		{
			name:       "SimpleAVL",
			tree:       simpleAVL,
			treeType:   treeTypeAVL,
			wantMin:    5,
			wantMax:    15,
			wantHeight: 2,
			wantNodes:  3,
			wantValues: true,
		},
		{
			name:       "SimpleRedBlack",
			tree:       simpleRedBlack,
			treeType:   treeTypeRedBlack,
			wantMin:    5,
			wantMax:    15,
			wantHeight: 2,
			wantNodes:  3,
			wantValues: true,
		},
		{
			name:       "DeepBST",
			tree:       deepBST,
			treeType:   treeTypeBST,
			wantMin:    3,
			wantMax:    20,
			wantHeight: 3,
			wantNodes:  6,
			wantValues: true,
		},
		{
			name:       "DeepAVL",
			tree:       deepAVL,
			treeType:   treeTypeAVL,
			wantMin:    3,
			wantMax:    20,
			wantHeight: 3,
			wantNodes:  6,
			wantValues: true,
		},
		{
			name:       "DeepRedBlack",
			tree:       deepRedBlack,
			treeType:   treeTypeRedBlack,
			wantMin:    3,
			wantMax:    20,
			wantHeight: 3,
			wantNodes:  6,
			wantValues: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := &Summary[int]{
				treeType:         "",
				nodeCount:        0,
				height:           tt.tree.Height(),
				isEmpty:          false,
				minValue:         0,
				maxValue:         0,
				hasValues:        false,
				balanceQuality:   BalanceUnknown,
				balanceScore:     0,
				optimalHeight:    0,
				heightEfficiency: 0,
			}
			summary.treeType = determineTreeType(tt.tree)

			calculateBasicMetrics(tt.tree, summary)

			if summary.nodeCount != tt.wantNodes {
				t.Errorf("nodeCount = %d, want %d", summary.nodeCount, tt.wantNodes)
			}

			if summary.hasValues != tt.wantValues {
				t.Errorf("hasValues = %t, want %t", summary.hasValues, tt.wantValues)
			}

			if summary.minValue != tt.wantMin {
				t.Errorf("minValue = %d, want %d", summary.minValue, tt.wantMin)
			}
			if summary.maxValue != tt.wantMax {
				t.Errorf("maxValue = %d, want %d", summary.maxValue, tt.wantMax)
			}

			// Optimal height: ceil(log2(nodes+1))
			wantOptimalHeight := int(math.Ceil(math.Log2(float64(summary.nodeCount + 1))))
			if summary.optimalHeight != wantOptimalHeight {
				t.Errorf("optimalHeight = %d, want %d", summary.optimalHeight, wantOptimalHeight)
			}
		})
	}
}
