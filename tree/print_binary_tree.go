package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// MaxPrintableLevel defines the maximum depth (0-indexed) of a tree that will
// be rendered by PrintBinaryTreeASCII. Trees deeper than this will be
// truncated with an ellipsis (...). This limit exists because the width of the
// ASCII tree grows exponentially with its height, quickly becoming too wide
// for standard displays.
const MaxPrintableLevel = 6 // Renders a tree of height 7 (levels 0 through 6).

// PrintBinaryTreeASCII generates a string containing an ASCII art representation
// of a binary tree. It is a useful tool for debugging and visualizing the
// structure of a tree.
//
// The rendering is optimized for node values that are up to 11 characters wide.
// An optional label can be provided, which will be printed before the tree.
// If the tree's height exceeds MaxPrintableLevel, it will be truncated.
//
// label is an optional title to print above the tree.
// t is the binary tree to be printed.
// Returns a string containing the ASCII representation.
func PrintBinaryTreeASCII[T constraints.Ordered](label string, t BinaryTree[T]) string {
	var buf bytes.Buffer
	if len(label) > 0 {
		buf.WriteString(label + "\n\n")
	}

	if isTreeNil(t) {
		buf.WriteString("nil tree\n")
		return buf.String()
	}

	stats := analyzeTree(t)
	height := stats.height
	if height < 0 {
		buf.WriteString("empty tree\n")
		return buf.String()
	}

	// Determine the starting depth for rendering, truncating if necessary.
	depthFrom := height
	if height > MaxPrintableLevel {
		depthFrom = MaxPrintableLevel
	}

	indentOpts := indentOptsForNodeWidth(stats.widestValue)
	nodes := []BinaryTree[T]{t}

	// Render the tree level by level, from the root down.
	outputNodes(nodes, indentOpts, &buf, depthFrom)
	for depthFrom > 0 {
		nodes = generateLevelsNodes(nodes)
		outputLegs(nodes, indentOpts, &buf, depthFrom)

		depthFrom--
		outputNodes(nodes, indentOpts, &buf, depthFrom)
	}

	// Add an ellipsis if the tree was taller than the print limit.
	if height > MaxPrintableLevel+1 {
		buf.WriteString("...\n")
	}

	return buf.String()
}

// writeLeg is a small helper to write either a leg segment or padding,
// depending on whether a node exists at a given position.
func writeLeg[T constraints.Ordered](leg BinaryTree[T], legString, padString string, buf *bytes.Buffer) {
	if !isTreeNil(leg) {
		buf.WriteString(legString)
	} else {
		buf.WriteString(padString)
	}
}

// outputLegs renders the connector lines (legs) between one level of nodes
// and the level below it. It calculates the correct spacing to ensure the
// legs point to the correct child positions.
func outputLegs[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int) {
	opts := indentOptions[depthFrom]
	lastNode := lastNonNilNode(nodes)
	if lastNode == -1 {
		return
	}

	// Render each segment of the legs (for multi-line legs).
	for i := 0; i < opts.legDepth; i++ {
		buf.WriteString(prefixPad[:opts.prefixPadding])
		for j := 0; j <= lastNode; j += 2 {
			// Calculate padding and leg characters for the left child.
			legDepthPad := opts.legDepth - 1 - i
			leftLegStr := legPad[:legDepthPad] + leftLegs[i]
			writeLeg(nodes[j], leftLegStr, indentFull[:opts.legDepth], buf)

			// Render padding between the two legs of a single parent.
			buf.WriteString(shoulderPad[:opts.shoulderPadding])
			buf.WriteString(intraPad[:opts.nodeWidth])
			buf.WriteString(shoulderPad[:opts.shoulderPadding])

			// Calculate padding and leg characters for the right child.
			rightLegStr := rightLegs[i] + legPad[:legDepthPad]
			if j+1 <= lastNode {
				writeLeg(nodes[j+1], rightLegStr, indentFull[:opts.legDepth], buf)
			} else {
				buf.WriteString(indentFull[:opts.legDepth])
			}

			// Render padding between different subtrees.
			if j+1 < lastNode {
				buf.WriteString(interPad[:opts.interTreePadding])
			}
		}
		buf.WriteString("\n")
	}
}

// outputNodes renders the node values and their metadata for a single level
// of the tree. It also renders the horizontal bars that connect to the legs.
func outputNodes[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int) {
	lastNode := lastNonNilNode(nodes)
	if lastNode == -1 {
		return
	}

	// --- Render Node Values ---
	printNodeLine(nodes, indentOptions, buf, depthFrom, true)

	// --- Render Node Metadata (if any) ---
	if levelHasMetadata(nodes) {
		printNodeLine(nodes, indentOptions, buf, depthFrom, false)
	}
}

// printNodeLine is a helper to render a single line of output for a level,
// either for the node values or their metadata.
func printNodeLine[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int, isValueLine bool) {
	opts := indentOptions[depthFrom]
	parentOpts := indentOptions[depthFrom+1]
	lastNode := lastNonNilNode(nodes)

	buf.WriteString(prefixPad[:opts.prefixPadding])
	for j, n := range nodes {
		// Determine the content to print (value or metadata).
		var content string
		if !isTreeNil(n) {
			if isValueLine {
				content = fmt.Sprintf(nodeFmt, n.Value())
			} else {
				content = n.Metadata()
			}
		}

		// Calculate inter-node spacing based on parent level's layout.
		inter := (parentOpts.legDepth + parentOpts.shoulderPadding) - (opts.legDepth + opts.shoulderPadding)

		// Add leading leg/shoulder padding.
		if depthFrom != 0 || (j%2 == 1) {
			buf.WriteString(legPad[:opts.legDepth])
		}
		if !isTreeNil(n) && n.HasLeft() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			buf.WriteString(shoulderPad[:opts.shoulderPadding])
		}

		// Write the centered content.
		if !isTreeNil(n) {
			buf.WriteString(centerString(content, " ", opts.nodeWidth))
		} else {
			buf.WriteString(indentFull[:opts.nodeWidth])
		}

		// Add trailing shoulder/leg padding.
		if !isTreeNil(n) && n.HasRight() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			buf.WriteString(shoulderPad[:opts.shoulderPadding])
		}
		buf.WriteString(legPad[:opts.legDepth])

		if j >= lastNode {
			break
		}

		// Add padding between nodes.
		if j%2 == 0 { // Space between a left and right child.
			buf.WriteString(shoulderPad[:inter])
			buf.WriteString(intraPad[:opts.nodeWidth])
			buf.WriteString(shoulderPad[:inter])
		} else { // Space between two different subtrees.
			buf.WriteString(interPad[:opts.interTreePadding])
		}
	}
	buf.WriteString("\n")
}

// levelHasMetadata checks if any node in a given level has non-empty metadata.
// This is used to decide whether to print the metadata line for a level.
func levelHasMetadata[T constraints.Ordered](nodes []BinaryTree[T]) bool {
	for _, n := range nodes {
		if !isTreeNil(n) && n.Metadata() != "" {
			return true
		}
	}
	return false
}