package tree

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

const treeTypeUnknown = "Unknown"

// BalanceQuality provides a qualitative assessment of a tree's balance.
// The categories are based on the balance score, providing a human-readable
// description of how skewed the tree is.
type BalanceQuality int

const (
	BalanceUnknown BalanceQuality = iota
	SeverelyLeftHeavy
	ModeratelyLeftHeavy
	SlightlyLeftHeavy
	WellBalanced
	SlightlyRightHeavy
	ModeratelyRightHeavy
	SeverelyRightHeavy
	Degenerate // Represents a tree that has devolved into a linked list.
)

// String returns the string representation of a BalanceQuality enum value.
func (bq BalanceQuality) String() string {
	switch bq {
	case BalanceUnknown:
		return "Unknown"
	case SeverelyLeftHeavy:
		return "Severely Left-Heavy"
	case ModeratelyLeftHeavy:
		return "Moderately Left-Heavy"
	case SlightlyLeftHeavy:
		return "Slightly Left-Heavy"
	case WellBalanced:
		return "Well-Balanced"
	case SlightlyRightHeavy:
		return "Slightly Right-Heavy"
	case ModeratelyRightHeavy:
		return "Moderately Right-Heavy"
	case SeverelyRightHeavy:
		return "Severely Right-Heavy"
	case Degenerate:
		return "Degenerate (Linear)"
	default:
		return "Unknown"
	}
}

// Summary provides a comprehensive analysis of a tree's structural properties
// and statistics. It is designed to give a quick, high-level understanding of
// the tree's state, including its size, balance, and efficiency.
type Summary[T constraints.Ordered] struct {
	treeType         string
	nodeCount        int
	height           int
	isEmpty          bool
	minValue         T
	maxValue         T
	hasValues        bool
	balanceQuality   BalanceQuality
	balanceScore     float64
	optimalHeight    int
	heightEfficiency float64
}

// TreeType returns the identified type of the tree (e.g., "BST", "AVL").
func (s *Summary[T]) TreeType() string {
	return s.treeType
}

// NodeCount returns the total number of nodes in the tree.
func (s *Summary[T]) NodeCount() int {
	return s.nodeCount
}

// Height returns the measured height of the tree.
func (s *Summary[T]) Height() int {
	return s.height
}

// IsEmpty reports whether the tree contains any nodes.
func (s *Summary[T]) IsEmpty() bool {
	return s.isEmpty
}

// MinValue returns the smallest value stored in the tree and a boolean
// indicating if the value is valid (i.e., the tree is not empty).
func (s *Summary[T]) MinValue() (T, bool) {
	return s.minValue, s.hasValues
}

// MaxValue returns the largest value stored in the tree and a boolean
// indicating if the value is valid.
func (s *Summary[T]) MaxValue() (T, bool) {
	return s.maxValue, s.hasValues
}

// HasValues reports whether the tree contains any values.
func (s *Summary[T]) HasValues() bool {
	return s.hasValues
}

// BalanceQuality returns the qualitative assessment of the tree's balance.
func (s *Summary[T]) BalanceQuality() BalanceQuality {
	return s.balanceQuality
}

// BalanceScore returns a numerical score from -1.0 (perfectly left-skewed) to
// 1.0 (perfectly right-skewed), with 0 representing a perfectly balanced tree.
func (s *Summary[T]) BalanceScore() float64 {
	return s.balanceScore
}

// OptimalHeight returns the theoretical minimum possible height for a binary
// tree with the same number of nodes.
func (s *Summary[T]) OptimalHeight() int {
	return s.optimalHeight
}

// HeightEfficiency returns the ratio of the optimal height to the actual
// height, providing a score from 0.0 to 1.0 where 1.0 is most efficient.
func (s *Summary[T]) HeightEfficiency() float64 {
	return s.heightEfficiency
}

// String provides a formatted, human-readable summary of the tree's statistics.
func (s *Summary[T]) String() string {
	if s == nil {
		return "nil summary"
	}
	var sb strings.Builder
	sb.WriteString("=== Tree Summary ===\n")
	sb.WriteString(fmt.Sprintf("  Type: %s\n", s.treeType))
	sb.WriteString(fmt.Sprintf("  Is Empty: %t\n", s.isEmpty))
	if !s.isEmpty {
		sb.WriteString(fmt.Sprintf("  Node Count: %d\n", s.nodeCount))
		sb.WriteString(fmt.Sprintf("  Height: %d (Optimal: %d, Efficiency: %.2f%%)\n", s.height, s.optimalHeight, s.heightEfficiency*100))
		if s.hasValues {
			sb.WriteString(fmt.Sprintf("  Value Range: [%v, %v]\n", s.minValue, s.maxValue))
		}
		sb.WriteString(fmt.Sprintf("  Balance: %s (Score: %.2f)\n", s.balanceQuality, s.balanceScore))
	}
	return sb.String()
}

// determineTreeType uses type assertion to identify the concrete type of the tree.
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

// calculateBasicMetrics traverses the tree to compute fundamental metrics like
// node count, value range, and height efficiency.
func calculateBasicMetrics[T constraints.Ordered](tree Tree[T], summary *Summary[T]) {
	if isTreeNil(tree) || tree.Height() < 0 {
		return
	}

	valueChan := tree.Traverse(TraverseInOrder)
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

	if summary.nodeCount > 0 {
		summary.optimalHeight = int(math.Ceil(math.Log2(float64(summary.nodeCount + 1))))
		if summary.height > 0 {
			summary.heightEfficiency = float64(summary.optimalHeight) / float64(summary.height)
		}
	}
}