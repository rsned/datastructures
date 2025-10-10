package tree

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

const treeTypeUnknown = "Unknown"

// BalanceQuality is a enum for the balance quality of a tree.
// It is an arbitrary set of values ranging from Well Balanced to
// Severely Right/Left Heavy. The change over points are chosen by
// me to what feels reasonable.
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
	Degenerate // Essentially a linked list
)

// balanceQualityString converts BalanceQuality enum to string
func balanceQualityString(bq BalanceQuality) string {
	switch bq {
	case BalanceUnknown:
		return treeTypeUnknown
	case SeverelyLeftHeavy:
		return "Severely Left Heavy"
	case ModeratelyLeftHeavy:
		return "Moderately Left Heavy"
	case SlightlyLeftHeavy:
		return "Slightly Left Heavy"
	case WellBalanced:
		return "Well Balanced"
	case SlightlyRightHeavy:
		return "Slightly Right Heavy"
	case ModeratelyRightHeavy:
		return "Moderately Right Heavy"
	case SeverelyRightHeavy:
		return "Severely Right Heavy"
	case Degenerate:
		return "Degenerate (Linear)"
	default:
		return treeTypeUnknown
	}
}

// Summary contains comprehensive statistics and analysis of a Tree's
// data structure to help developers understand tree characteristics at a
// glance.
type Summary[T constraints.Ordered] struct {
	// Basic Tree Information
	treeType  string // "BST", "AVL", "RedBlack", etc.
	nodeCount int    // Total number of nodes
	height    int    // Maximum depth from root to leaf
	isEmpty   bool   // Whether tree has any nodes

	// TODO(rsned): Add more fields of interest.
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

// String returns a human-readable representation of the tree summary
func (ts *Summary[T]) String() string {
	var sb strings.Builder

	sb.WriteString("=== Tree Summary ===\n")
	sb.WriteString(fmt.Sprintf("Type: %s\n", ts.treeType))
	sb.WriteString(fmt.Sprintf("Empty: %t\n", ts.isEmpty))

	if !ts.isEmpty {
		sb.WriteString(fmt.Sprintf("Nodes: %d\n", ts.nodeCount))
		sb.WriteString(fmt.Sprintf("Height: %d\n", ts.height))
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
