package tree

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

const treeTypeUnknown = "Unknown"

// Summary contains comprehensive statistics and analysis of a Tree's
// data structure to help understand its characteristics at a glance.
type Summary[T constraints.Ordered] struct {
	// Basic Tree Information
	treeType  string // "BST", "AVL", "RedBlack", etc.
	nodeCount int    // Total number of nodes
	height    int    // Maximum depth from root to leaf
	isEmpty   bool   // Whether tree has any nodes

	// Value Range Information
	minValue  T    // Smallest value in tree
	maxValue  T    // Largest value in tree
	hasValues bool // Whether min/max are valid

	// Balance and Structure Metrics
	balanceQuality   BalanceQuality // Enum: Balanced, LeftHeavy, RightHeavy, etc.
	balanceScore     float64        // Numerical balance score [-1.0 to 1.0]
	optimalHeight    int            // Theoretical minimum height for node count
	heightEfficiency float64        // Ratio of optimal to actual height [0 to 1.0]

	// TODO(rsned): Add more fields of interest.
}

// zeroValue returns the zero value for the given type
func zeroValue[T constraints.Ordered]() T {
	return *new(T)
}

// newSummary creates a new instance of a Summary ready to be filled.
func newSummary[T constraints.Ordered]() *Summary[T] {
	return &Summary[T]{
		treeType:         treeTypeUnknown,
		nodeCount:        0,
		height:           0,
		isEmpty:          true,
		minValue:         zeroValue[T](),
		maxValue:         zeroValue[T](),
		hasValues:        false,
		balanceQuality:   BalanceUnknown,
		balanceScore:     0.0,
		optimalHeight:    0,
		heightEfficiency: 0.0,
	}
}

// TreeType returns the type of tree (BST, AVL, RedBlack, etc.)
func (ts *Summary[T]) TreeType() string {
	return ts.treeType
}

// NodeCount returns the total number of nodes in the tree
func (ts *Summary[T]) NodeCount() int {
	return ts.nodeCount
}

// Height returns the maximum depth from root to leaf
func (ts *Summary[T]) Height() int {
	return ts.height
}

// IsEmpty returns whether the tree has any nodes
func (ts *Summary[T]) IsEmpty() bool {
	return ts.isEmpty
}

// MinValue returns the minimum value in the tree
func (ts *Summary[T]) MinValue() (T, bool) {
	return ts.minValue, ts.hasValues
}

// MaxValue returns the maximum value in the tree
func (ts *Summary[T]) MaxValue() (T, bool) {
	return ts.maxValue, ts.hasValues
}

// HasValues returns whether the tree contains values
func (ts *Summary[T]) HasValues() bool {
	return ts.hasValues
}

// BalanceQuality returns the balance quality enum
func (ts *Summary[T]) BalanceQuality() BalanceQuality {
	return ts.balanceQuality
}

// BalanceScore returns the numerical balance score (-1.0 to 1.0)
func (ts *Summary[T]) BalanceScore() float64 {
	return ts.balanceScore
}

// OptimalHeight returns the theoretical minimum height for the node count
func (ts *Summary[T]) OptimalHeight() int {
	return ts.optimalHeight
}

// HeightEfficiency returns the ratio of optimal to actual height
func (ts *Summary[T]) HeightEfficiency() float64 {
	return ts.heightEfficiency
}

// String returns a human-readable representation of the tree summary
func (ts *Summary[T]) String() string {
	var sb strings.Builder

	sb.WriteString("=== Tree Summary ===\n")
	sb.WriteString(fmt.Sprintf("Type: %s\n", ts.treeType))
	sb.WriteString(fmt.Sprintf("Empty: %t\n", ts.isEmpty))

	if !ts.isEmpty {
		sb.WriteString(fmt.Sprintf("Nodes: %d\n", ts.nodeCount))
		sb.WriteString(fmt.Sprintf("Height: %d (optimal: %d)\n", ts.height, ts.optimalHeight))
		sb.WriteString(fmt.Sprintf("Height Efficiency: %.1f%%\n", ts.heightEfficiency*100))

		if ts.hasValues {
			sb.WriteString(fmt.Sprintf("Value Range: %v to %v\n", ts.minValue, ts.maxValue))
		}

		sb.WriteString("\n=== Balance Analysis ===\n")
		sb.WriteString(fmt.Sprintf("Balance Quality: %s\n", ts.balanceQuality.String()))
		sb.WriteString(fmt.Sprintf("Balance Score: %.3f\n", ts.balanceScore))
		sb.WriteString(BalanceQualityGraph(ts.balanceScore) + "\n")
	}

	return sb.String()
}

// determineTreeType identifies the concrete type of the tree
func determineTreeType[T constraints.Ordered](tree Tree[T]) string {
	switch tree.(type) {
	case *AVL[T]:
		return treeTypeAVL
	case *BST[T]:
		return treeTypeBST
	case *RedBlack[T]:
		return treeTypeRedBlack
	default:
		return treeTypeUnknown
	}
}

// calculateBasicMetrics computes node count and value range
//
// This method is not thread-safe as the Tree could be modified during the
// traversing of the tree leading to inconsistent results.
func calculateBasicMetrics[T constraints.Ordered](tree Tree[T], summary *Summary[T]) {
	summary.height = tree.Height()
	// Check if tree is empty first (avoids interface nil pointer issues)
	if tree.Height() == 0 {
		return
	}

	// Get traverser interface - all Tree implementations provide this
	traverser, ok := tree.(Traverser[T])
	if !ok {
		return
	}

	// Use the tree's built-in traversal (order doesn't matter for metrics)
	valueChan := traverser.Traverse(TraverseInOrder)

	// Read from channel and update min/max values
	first := true
	for v := range valueChan {
		summary.nodeCount++
		if first {
			summary.hasValues = true
			summary.minValue = v
			summary.maxValue = v
			first = false
		} else {
			if v < summary.minValue {
				summary.minValue = v
			}
			if v > summary.maxValue {
				summary.maxValue = v
			}
		}
	}

	// Calculate optimal height
	if summary.nodeCount > 0 {
		summary.optimalHeight = int(math.Ceil(math.Log2(float64(summary.nodeCount + 1))))
		if summary.height > 0 {
			summary.heightEfficiency = float64(summary.optimalHeight) / float64(summary.height)
		}
	}
	if summary.height > 0 {
		summary.isEmpty = false
	}
}
