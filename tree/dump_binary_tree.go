package tree

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// This file contains the logic for rendering a visual representation of a
// binary tree as an ASCII string. It is primarily used for debugging and
// educational purposes to visualize the structure of a tree.

const (
	// Basic string components for building the tree structure.
	indent      = "     "
	nodeFmt     = "%-3v"
	nodeMetaFmt = "%5s"
	underbar    = "_____"

	// Leg segments for drawing connector lines between nodes.
	leftRow1  = "/"
	leftRow2  = "/ "
	leftRow3  = "/  "
	rightRow1 = "\\"
	rightRow2 = " \\"
	rightRow3 = "  \\"
)

// maxPadding defines the pre-allocated size for repeated string patterns
// to optimize rendering performance by avoiding repeated string concatenation.
const maxPadding = 2048

var (
	// Slices of leg strings for different depths, allowing for easy iteration.
	leftLegs  = []string{leftRow1, leftRow2, leftRow3}
	rightLegs = []string{rightRow1, rightRow2, rightRow3}

	// Pre-built strings for padding, used via substring operations for efficiency.
	underbarFull = strings.Repeat(underbar, maxPadding)
	indentFull   = strings.Repeat(indent, maxPadding)
	prefixPad    = strings.Repeat(" ", maxPadding)
	shoulderPad  = strings.Repeat(" ", maxPadding)
	interPad     = strings.Repeat(" ", maxPadding)
	intraPad     = strings.Repeat(" ", maxPadding)
	legPad       = strings.Repeat(" ", maxPadding)
)

// indentOptions defines the spacing and padding parameters used at a given
// depth of the tree to ensure proper alignment and a visually clear layout.
// The fields correspond to different padding areas around and between nodes.
type indentOptions struct {
	nodeWidth        int // Width of the node's value representation.
	prefixPadding    int // Initial padding for the entire line.
	intraNodePadding int // Padding between the legs of a single node.
	interTreePadding int // Padding between different subtrees at the same level.
	shoulderPadding  int // Horizontal padding to align nodes under their parents.
	legDepth         int // Vertical height of the connector legs.
}

// RenderMode specifies the output format for tree visualization.
type RenderMode int

const (
	ModeASCII RenderMode = iota // Default mode, renders as plain text ASCII art.
	ModeSVG                     // Renders as a Scalable Vector Graphic (not implemented).
)

// RenderBinaryTree generates a string representation of a binary tree in the
// specified mode. Currently, only ASCII mode is implemented.
func RenderBinaryTree[T constraints.Ordered](t BinaryTree[T], _ int, mode RenderMode) string {
	switch mode {
	case ModeASCII:
		return PrintBinaryTreeASCII("", t)
	case ModeSVG:
		return "SVG rendering is not yet implemented."
	default:
		return "Unknown render mode."
	}
}

// maxNodeWidth sets the maximum supported width for a node's value string.
// Values wider than this may not render correctly.
const maxNodeWidth = 11

// indentOptsForNodeWidth retrieves the appropriate indentation and spacing
// configuration for a given node value width from a pre-computed data table.
func indentOptsForNodeWidth(width int) indentOptionsMap {
	if width < 1 {
		width = 1
	} else if width > maxNodeWidth {
		width = maxNodeWidth
	}
	// The spacing data is indexed by width, with adjustments for odd/even widths.
	return binaryTreeSpacingData[width+(width+1)%2]
}

// generateLevelsNodes creates a slice representing the next level of the tree
// based on the nodes in the current level. It preserves the tree structure by
// using nil placeholders for missing children, resulting in a sparse slice.
func generateLevelsNodes[T constraints.Ordered](existing []BinaryTree[T]) []BinaryTree[T] {
	nodes := make([]BinaryTree[T], 0, len(existing)*2)
	for _, n := range existing {
		if !isTreeNil(n) {
			nodes = append(nodes, n.Left(), n.Right())
		} else {
			// Add placeholders for children of a nil node to maintain spacing.
			nodes = append(nodes, nil, nil)
		}
	}
	return nodes
}

// lastNonNilNode finds the index of the rightmost non-nil node in a level.
// This is an optimization to avoid rendering unnecessary trailing spaces.
func lastNonNilNode[T constraints.Ordered](nodes []BinaryTree[T]) int {
	for i := len(nodes) - 1; i >= 0; i-- {
		if !isTreeNil(nodes[i]) {
			return i
		}
	}
	return -1 // Indicates an entirely nil level.
}

// centerString centers a string `s` within a given `width` by padding it
// with the specified `padChar`. If the string is longer than the width, it
// is returned unchanged.
func centerString(s, padChar string, width int) string {
	s = strings.TrimSpace(s)
	l := len(s)
	if l >= width {
		return s
	}
	lPad := (width - l) / 2
	rPad := width - l - lPad
	return strings.Repeat(padChar, lPad) + s + strings.Repeat(padChar, rPad)
}

// dumpTreeStats holds metrics about a tree's structure, used to guide rendering.
type dumpTreeStats struct {
	height      int
	leftHeight  int
	rightHeight int
	widestValue int
}

// analyzeTree traverses a tree to gather statistics like height and the width
// of the widest node value. These stats help in formatting the ASCII output.
func analyzeTree[T constraints.Ordered](tree BinaryTree[T]) dumpTreeStats {
	if isTreeNil(tree) {
		return dumpTreeStats{}
	}

	stats := dumpTreeStats{height: tree.Height()}
	if tree.HasLeft() {
		stats.leftHeight = tree.Left().Height()
	}
	if tree.HasRight() {
		stats.rightHeight = tree.Right().Height()
	}

	var widest int
	ch := tree.Traverse(TraverseLevelOrder) // Level order is efficient for this.
	for val := range ch {
		s := fmt.Sprintf("%v", val)
		if len(s) > widest {
			widest = len(s)
		}
	}
	stats.widestValue = widest
	return stats
}