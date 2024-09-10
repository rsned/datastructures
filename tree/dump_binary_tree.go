package tree

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// TODO(rsned): A list of potential enhancements.
//
// * Find widest node value to be able to change the overall height and width of
//   the output tree. Shorter values/strings don't need as tall or wide of a tree.
// * Find the breadth of a given subtree and use it to adjust the lateral width
//   of higher up nodes.  e.g. when one side of a tree is not bushy, or is
//   unbalanced, there is no need for lateral padding on higher nodes.
// * Allow for pseudo-dynamic heights based on width of largest element in the tree.
//   e.g. if the tree only has single letter / digit values, a leg height of 2-3
//   would be plenty.
// * Format a nodes content to be centered in the alloted space rather than left
//   or right aligned.
// * Node value and metadata printing are basically identical code blocks, figure
//   out a way to refactor that.

const (
	// How wide is the unit of ascii art we are using.
	elementWidth = 5

	indent      = "     "
	nodeFmt     = "~%3d~"
	nodeFmtT    = "%3v"
	nodeMetaFmt = "%5s"

	leftLegBase  = "/"
	rightLegBase = "\\"

	leftRow1  = "/"
	leftRow2  = "/ "
	leftRow3  = "/  "
	leftRow4  = "/   "
	leftRow5  = "/    "
	rightRow1 = "\\"
	rightRow2 = " \\"
	rightRow3 = "  \\"
	rightRow4 = "   \\"
	rightRow5 = "    \\"
	rightRow6 = "     \\"
	rightRow7 = "      \\"
	underbar  = "_____"
)

var (
	// leftLegs is a slice of the angled leg strings in order to be
	// iterated over at each level to make the render function cleaner.
	leftLegs = []string{
		leftRow1,
		leftRow2,
		leftRow3,
		leftRow4,
		leftRow5,
	}

	// the corresponding right angled leg strings.
	rightLegs = []string{
		rightRow1,
		rightRow2,
		rightRow3,
		rightRow4,
		rightRow5,
		rightRow6,
		rightRow7,
	}

	// Thise are constructed to allow substring instead of looping repeatedly
	// when multiple instances are needed in a row.
	underbarFull = strings.Repeat(underbar, 40)
	indentFull   = strings.Repeat(indent, 40)
	prefixPad    = strings.Repeat("P", 40)
	shoulderPad  = strings.Repeat("S", 40)
	interPad     = strings.Repeat("I", 40)
	intraPad     = strings.Repeat("i", 40)
	otherPad     = strings.Repeat("#", 40)
	otherPad2    = strings.Repeat("$", 40)
	legPad       = strings.Repeat("L", 40)
)

// indentOptions tracks the spacings used at a given depth and tree height for a given node width.
type indentOptions struct {
	// indentWidth is how wide in number of spaces one indent "unit" is for this
	// set of options.  For trees with all narrow values, the width will be smaller.
	indentWidth int

	// prefixPadding is how much spacing to start the beginning of a line with.
	//
	// This is measured in spaces.
	prefixPadding int

	// intraNodePadding is the spacing between the left and right legs of the tree.
	//
	// This is measured in spaces.
	intraNodePadding int

	// interTreePadding is the spacing between each set of trees at this level.
	//
	// This is measured in units of indentWidth
	interTreePadding int

	// shoulderPadding is how much lateral filler we need between the top of a leg
	// and the levels node. (rather than growing the diagonal 2^n vertically, we
	// limit it to 5 vertical levels and then go sideways to make up the space.)
	//
	// This is in spaces.
	shoulderPadding int

	// legDepth tracks how many of the vertical levels of the legs to show.
	legDepth int
}

const (
	depthBottom      = 0
	depthBottomPlus1 = 1
	depthBottomPlus2 = 2
	depthBottomPlus3 = 3
	depthBottomPlus4 = 4
)

// RenderMode is an enum for potential outputs when dumping or rendering the trees.
type RenderMode int

// Set of current render modes.
const (
	ModeASCII RenderMode = iota
	ModeSVG

	// TODO(rsned): Add more modes?
)

// RenderBinaryTree returns the given tree in the given mode rendered into string form.
func RenderBinaryTree[T constraints.Ordered](t BinaryTree[T], height int, mode RenderMode) string {
	switch mode {
	case ModeASCII:
		return dumpBinaryTree("", t)
	default:
		return "Method not implemented yet"
	}
}

// dumpBinaryTree is a simple hacky way to output a binary tree up to 5 levels
// for the purpose of aiding in testing and debugging.
//
// This method takes a binary tree of integers with an explicit assumption the
// values in the tree are 3 character or less wide when printed.
// e.g. -21, 123, 7, etc.
//
// An optional label is output before the tree contents.
func dumpBinaryTree[T constraints.Ordered](label string, t BinaryTree[T]) string {
	var buf bytes.Buffer
	// This doesn't work on interface to generic types.
	// If the tree is nil, it skips this and crashes later on.
	if isTreeNil(t) {
		return buf.String()
	}

	stats := analyzeTree(t)
	height := stats.height
	node := t
	nodes := []BinaryTree[T]{node}
	var nextNodes []BinaryTree[T]
	depthFrom := height - 1
	indentOpts := optsForStats(depthFrom, stats.widestValue)

	// First pass starts with the root node, then we go into the loop of
	// legs and nodes until we are all done.
	outputNodes(nodes, indentOpts, &buf, depthFrom)

	for depthFrom > 0 {
		nextNodes = generateLevelsNodes(nodes)
		outputLegs(nextNodes, indentOpts, &buf, depthFrom)

		depthFrom--
		indentOpts = optsForStats(depthFrom, stats.widestValue)
		nodes = nextNodes
		outputNodes(nodes, indentOpts, &buf, depthFrom)
	}

	return buf.String()
}

func optsForStats(depthFrom, widest int) indentOptionsMap {
	if widest <= 1 {
		return binaryTreeSpacingData[1]
	}
	if widest <= 3 {
		return binaryTreeSpacingData[3]
	}
	return binaryTreeSpacingData[5]
}

// generateLevelsNodes ranges over the given set of nodes generating a new
// set of nodes from the children. Nodes which don't exist are replaced with
// nils as placeholders for nodes that are not in the tree.
//
// This slice will be a potentially sparse slice of size 2*N
// where N is the number of nodes from the previous level.
func generateLevelsNodes[T constraints.Ordered](existing []BinaryTree[T]) []BinaryTree[T] {
	nodes := []BinaryTree[T]{}
	for _, n := range existing {
		if n != nil {
			if n.HasLeft() {
				nodes = append(nodes, n.Left())
			} else {
				nodes = append(nodes, nil)
			}
			if n.HasRight() {
				nodes = append(nodes, n.Right())
			} else {
				nodes = append(nodes, nil)
			}
		} else {
			// If the parent nodes had a gap (i.e. an unbalanced
			// tree in progress), we add two blank filler entries
			// in its place.
			nodes = append(nodes, nil)
			nodes = append(nodes, nil)
		}
	}

	return nodes
}

// lastNonNilNode walks backward through the list looking for the first non-nil
// entry in it so we can cut the row processing early if possible.
func lastNonNilNode[T constraints.Ordered](nodes []BinaryTree[T]) int {
	for i := len(nodes) - 1; i >= 0; i-- {
		if nodes[i] != nil {
			return i
		}
	}

	return 0
}

// writeLeg replaces the boilerplate with a simple helper.
func writeLeg[T constraints.Ordered](leg BinaryTree[T], legString string, indentString string, buf *bytes.Buffer) {
	if leg != nil {
		buf.WriteString(legString)
	} else {
		//buf.WriteString(indentString)
		buf.WriteString(legString)
	}
}

// outputLegs does the boring bits of printing out visible or missing legs and the
// appropriate spacings between each one.
func outputLegs[T constraints.Ordered](nodes []BinaryTree[T], indentOptions indentOptionsMap, buf *bytes.Buffer, depthFrom int) {
	opts := indentOptions[depthFrom]
	nodeSize := opts.indentWidth
	lastNode := lastNonNilNode(nodes)
	for i, ll := range leftLegs[:opts.legDepth] {
		buf.WriteString(prefixPad[:opts.prefixPadding])
		for j := 0; j < len(nodes); j++ {
			if j > lastNode {
				break
			}

			legDepthPad := opts.legDepth - 1 - i

			// offset is based on number of legs to be drawn at this level.
			// left leg needs to be limited to this legDepth.
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
	nodeSize := opts.indentWidth
	parentOpts := indentOptions[depthFrom+1]
	lastNode := lastNonNilNode(nodes)

	// Nodes.
	buf.WriteString(prefixPad[:opts.prefixPadding])
	for j, n := range nodes {
		// For all rows except the bottom row,  each node potentially has
		// both left and right legs below it that need to be padded for.
		if depthFrom != 0 || (depthFrom == 0 && j != 0 && j%2 == 1) {
			buf.WriteString(legPad[:opts.legDepth])
			//} else {
			//buf.WriteString("*")
		}

		// Higher up levels have lines that go sideways to keep the tree
		// reasonably sized.
		if n != nil && n.HasLeft() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			//buf.WriteString(shoulderPad[:opts.shoulderPadding])
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		}

		// The actual node value.
		if n != nil {
			buf.WriteString(centerString(fmt.Sprintf(nodeFmtT, n.Value()), " ",
				nodeSize))
		} else {
			buf.WriteString(indentFull[:nodeSize])
		}

		if n != nil && n.HasRight() {
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		} else {
			//buf.WriteString(shoulderPad[:opts.shoulderPadding])
			buf.WriteString(underbarFull[:opts.shoulderPadding])
		}
		// If this is the last node, skip all the remaining trailing padding.
		if j >= lastNode {
			break
		}

		// This is the padding to match the leg above it.
		// If this is an even index, then we want the padding to match
		// the number of leg segments leading down into this node
		// on the inside of the node values.
		// if j%2 == 0 {
		buf.WriteString(legPad[:opts.legDepth])
		// } else {
		// buf.WriteString("*")
		// }

		// Between the even and odd indexes the spacing breakdown
		// matches what the outputLegs does (combination of shoulder
		// spacing and nodeWidth but based on the next higher level
		// ups indent optiond. e.g Even index values represent
		// left legs and odd indexes represent right legs.
		//
		inter := (parentOpts.legDepth + parentOpts.shoulderPadding) -
			(opts.legDepth + opts.shoulderPadding)
		if j%2 == 0 {
			//buf.WriteString(shoulderPad[:opts.shoulderPadding])
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

	// Add metadata print
	buf.WriteString(prefixPad[:opts.prefixPadding])
	for j, n := range nodes {
		// This indent lines up with the left leg lines.
		if depthFrom != 0 {
			buf.WriteString(indentFull[:nodeSize])
			// buf.WriteString(indent)
		}
		buf.WriteString(shoulderPad[:opts.shoulderPadding])
		if n != nil {
			buf.WriteString(fmt.Sprintf(nodeMetaFmt, n.Metadata()))
		} else {
			buf.WriteString(indentFull[:nodeSize])
			// buf.WriteString(indent)
		}
		// If this is the last node, skip all the remaining trailing padding.
		if j >= lastNode {
			break
		}

		buf.WriteString(shoulderPad[:opts.shoulderPadding])
		// This indent lines up with the right leg lines.
		if depthFrom != 0 {
			buf.WriteString(indentFull[:nodeSize])
			// buf.WriteString(indent)
		}
		buf.WriteString(interPad[:opts.interTreePadding])
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

// centerString centers the given string into the target size adjusting the space at
// either end as needed with the given pad character.
//
// This method assumes an output width of 50 or less for the purpose of this file.
func centerString(s, padChar string, width int) string {
	s = strings.TrimSpace(s)
	l := len(s)

	// For now, there is no attempt to truncate or elide longer values.
	if l >= width {
		return s
	}

	// We attempt to right justify uneven splits so values will sit
	// slightly more right than left when they can't be balanced.
	lPad := ((width - l) / 2) + 1
	rPad := width - l - lPad

	// TODO(rsned): replace spaces with the padChar
	const spaces = "                                                                     "
	return fmt.Sprintf("%s%s%s", spaces[0:lPad], s, spaces[0:rPad])
}

type dumpTreeStats struct {
	height      int
	leftHeight  int
	rightHeight int
	widestValue int
}

// analyzeTree takes the givern tree and attempts to find out relevant details
// about it to assist in the rendering.
func analyzeTree[T constraints.Ordered](tree BinaryTree[T]) dumpTreeStats {

	stats := dumpTreeStats{
		height:      tree.Height(),
		leftHeight:  tree.Left().Height(),
		rightHeight: tree.Right().Height(),
	}

	// things we want to find out:
	// max height
	// width of largest value
	// lopsidedness / skew   e.g. is this only a one sided binary tree?

	// Walk the tree getting all values and printing them as strings.
	ch := tree.Traverse(TraverseInOrder)

	var widest int
	var got []string
	for {
		val, ok := <-ch
		if ok {
			s := fmt.Sprintf("%v", val)
			if len(s) > widest {
				widest = len(s)
			}
			got = append(got, fmt.Sprintf("%v", s))
		} else {
			break
		}
	}

	stats.widestValue = widest

	return stats
}
