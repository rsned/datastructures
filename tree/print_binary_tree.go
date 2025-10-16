package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// MaxPrintableLevel is the maximum number of levels to print because
// the ASCII tree is approximately doubling in printed width at each level.
//
// This value is arbitrary.
//
// But using the general printing spacings for a node width of various chars,
// we end up with a potential max width for depth N as follows.
//
// depth | width=3 | width=5 | width=7
// ------+---------+---------+---------
// 1     |       3 |       5 |       7
// 2     |      12 |      13 |      24
// 3     |      26 |      33 |      50
// 4     |      54 |      73 |     102
// 5     |     110 |     158 |     206
// 6     |     222 |     318 |     414
// 7     |     443 |     638 |     822
// 8     |     894 |    1278 |    1662
// 9     |    1790 |    2558 |    3326
const MaxPrintableLevel = 6 // 0-based, so 6 means 7 levels of tree height (0-6)

// PrintBinaryTreeASCII outputs the binary tree up to 6 levels deep
// for the purpose of aiding in testing and debugging.
//
// This method takes a binary tree of type T with an explicit assumption
// the values in the tree are under 11 character wide when formatted and
// printed. e.g. -21, 123, 7, 6.02e23, etc.  Values wider than 11 characters
// may work fine, but no testing or quality checks have been done.
//
// An optional label is output before the tree contents with a blank line
// separator between the label and the tree.
//
// TODO(rsned): Add an elideLevel int param to do the eliding dynamically.
func PrintBinaryTreeASCII[T constraints.Ordered](label string, t BinaryTree[T]) string {
	var buf bytes.Buffer
	if len(label) > 0 {
		buf.WriteString(label + "\n\n")
	}
	// This doesn't work on interface to generic types.
	// If the tree is nil, it skips this and crashes later on.
	if isTreeNil(t) {
		buf.WriteString("nil tree\n")

		return buf.String()
	}

	stats := analyzeTree(t)
	height := stats.height

	node := t
	nodes := []BinaryTree[T]{node}
	var nextNodes []BinaryTree[T]

	var depthFrom int
	if height > MaxPrintableLevel {
		depthFrom = MaxPrintableLevel
	} else {
		depthFrom = height - 1
	}

	indentOpts := indentOptsForNodeWidth(stats.widestValue)

	// First pass starts with the root node, then we go into the loop of
	// legs and nodes until we are all done.
	outputNodes(nodes, indentOpts, &buf, depthFrom)

	for depthFrom > 0 {
		nextNodes = generateLevelsNodes(nodes)
		outputLegs(nextNodes, indentOpts, &buf, depthFrom)

		depthFrom--
		nodes = nextNodes
		outputNodes(nodes, indentOpts, &buf, depthFrom)
	}

	// Add ellipsis if tree was truncated
	if height > MaxPrintableLevel+1 {
		buf.WriteString("...\n")
	}

	return buf.String()
}

// writeLeg replaces the boilerplate with a simple helper.
func writeLeg[T constraints.Ordered](leg BinaryTree[T], legString string, indentString string, buf *bytes.Buffer) {
	if leg != nil {
		buf.WriteString(legString)
	} else {
		buf.WriteString(indentString)
	}
}

// outputLegs does the boring bits of printing out visible or missing
// legs and the appropriate spacings between each one.
func outputLegs[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int) {
	opts := indentOptions[depthFrom]
	nodeSize := opts.nodeWidth
	lastNode := lastNonNilNode(nodes) // Tells us when to stop printing a line.

	for i, ll := range leftLegs[:opts.legDepth] {
		buf.WriteString(prefixPad[:opts.prefixPadding])
		for j := 0; j < len(nodes); j++ {
			// If we've handled the last real node in the list, break out.
			if j > lastNode {
				break
			}

			legDepthPad := opts.legDepth - 1 - i

			// offset is based on number of leg segments to be drawn at
			// this level. left leg needs to be limited to this legDepth.
			leftLeg := otherPad[:legDepthPad] + ll
			writeLeg(nodes[j], leftLeg, indentFull[:opts.legDepth], buf)

			// If this level has lateral legs, put in blanks to cover.
			buf.WriteString(shoulderPad[:opts.shoulderPadding])

			// Right legs are the next value, so jump forward to them.
			j++
			if j > lastNode {
				break
			}
			// Double check that we don't have an odd number of nodes.
			if j >= len(nodes) {
				break
			}

			// The spacing between the two legs in the tree.
			// Higher up nodes in the tree have more spacing to handle
			// the fanout as the tree grows.
			buf.WriteString(intraPad[:nodeSize])

			// If this level has lateral leg elements, put in blanks to cover.
			buf.WriteString(shoulderPad[:opts.shoulderPadding])

			// right leg needs to be limited to legDepth
			rl := rightLegs[i] + otherPad2[:legDepthPad]
			writeLeg(nodes[j], rl, indentFull[:opts.legDepth], buf)

			// For all but the final node in the list.
			if j != len(nodes)-1 {
				// Spacing between subtrees.
				buf.WriteString(interPad[:opts.interTreePadding])
			}
		}
		buf.WriteString("\n")
	}
}

// outputNodes writes out all the nodes and metadata at this level.
func outputNodes[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int) {
	opts := indentOptions[depthFrom]
	nodeSize := opts.nodeWidth
	parentOpts := indentOptions[depthFrom+1]
	lastNode := lastNonNilNode(nodes) // Tells us when to stop printing a line.

	// TODO(rsned): printing of nodes and printing of metadata is a giant
	// copy and paste.  Replace with one block of logic. Pre-generate a
	// slice of node formatted values and a slice of metadata formatted
	// values and then run both slices through the same print loop.
	buf.WriteString(prefixPad[:opts.prefixPadding])
	for j, n := range nodes {
		// For all rows except the bottom row, each node potentially has
		// both left and right legs below it that need to be padded for.
		if depthFrom != 0 || (depthFrom == 0 && j != 0 && j%2 == 1) {
			buf.WriteString(legPad[:opts.legDepth])
		}

		// Levels more than 3 from the bottom have "leg" lines that go
		// sideways to keep the tree reasonably compact vertically.
		if n != nil && n.HasLeft() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			buf.WriteString(shoulderPad[:opts.shoulderPadding])
		}

		// The actual node value.
		if n != nil {
			buf.WriteString(centerString(fmt.Sprintf(nodeFmt, n.Value()), " ",
				nodeSize))
		} else {
			buf.WriteString(indentFull[:nodeSize])
		}

		// If this is the last node of the line with no right child,
		// skip all the remaining work.
		if j >= lastNode && n != nil && !n.HasRight() {
			break
		}

		// This is the padding or right child underbar.
		if n != nil && n.HasRight() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			buf.WriteString(shoulderPad[:opts.shoulderPadding])
		}

		// We want the padding to match the number of leg segments
		// leading down into the child nodes.
		buf.WriteString(legPad[:opts.legDepth])

		// If we've handled the last real node in the list, break out without
		// add more trailing padded we don't need.
		if j >= lastNode {
			break
		}

		// Between the even and odd node indexes the spacing breakdown
		// matches what the outputLegs does (combination of shoulder
		// spacing and nodeWidth but based on the next higher level
		// ups indent optiond. e.g Even index values represent
		// left legs and odd indexes represent right legs.
		inter := (parentOpts.legDepth + parentOpts.shoulderPadding) -
			(opts.legDepth + opts.shoulderPadding)
		if j%2 == 0 {
			buf.WriteString(shoulderPad[:inter])
			buf.WriteString(intraPad[:nodeSize])
			buf.WriteString(shoulderPad[:inter])
		} else {
			// Finish off with the spacing between the trees.
			buf.WriteString(interPad[:opts.interTreePadding])
		}
	}
	buf.WriteString("\n")

	if !levelHasMetadata(nodes) {
		return
	}

	// TODO(rsned): Need to convert the print of metadata to match the spacing / alignment as the values.

	// Add metadata print
	buf.WriteString(prefixPad[:opts.prefixPadding])
	for j, n := range nodes {
		// For all rows except the bottom row,  each node potentially has
		// both left and right legs below it that need to be padded for.
		if depthFrom != 0 || (depthFrom == 0 && j != 0 && j%2 == 1) {
			buf.WriteString(legPad[:opts.legDepth])
		}

		// Higher up levels have lines that go sideways to keep the tree
		// reasonably sized.
		buf.WriteString(shoulderPad[:opts.shoulderPadding])

		if n != nil {
			buf.WriteString(centerString(fmt.Sprintf(nodeMetaFmt, n.Metadata()), " ", nodeSize))
		} else {
			buf.WriteString(indentFull[:nodeSize])
		}

		// If this is the last node, skip all the remaining trailing padding.
		if j >= lastNode {
			break
		}

		// Higher up levels have lines that go sideways to keep the tree
		// reasonably sized.
		buf.WriteString(shoulderPad[:opts.shoulderPadding])

		// This is the padding to match the leg above it.
		// If this is an even index, then we want the padding to match
		// the number of leg segments leading down into this node
		// on the inside of the node values.
		buf.WriteString(legPad[:opts.legDepth])

		if j >= lastNode {
			break
		}

		inter := (parentOpts.legDepth + parentOpts.shoulderPadding) -
			(opts.legDepth + opts.shoulderPadding)
		if j%2 == 0 {
			buf.WriteString(shoulderPad[:inter])
			buf.WriteString(intraPad[:nodeSize])
			buf.WriteString(shoulderPad[:inter])
		} else {
			// Finish off with the spacing between the trees.
			buf.WriteString(interPad[:opts.interTreePadding])
		}
	}
	buf.WriteString("\n")
}

// levelHasMetadata reports if the current set of nodes has any elements with
// some metadata value.
func levelHasMetadata[T constraints.Ordered](nodes []BinaryTree[T]) bool {
	has := false
	for _, n := range nodes {
		if n == nil {
			continue
		}
		has = has || (n.Metadata() != "")
	}

	return has
}
