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
	nodeFmt     = " %3d "
	nodeFmtT    = " %3v "
	nodeMetaFmt = "%5s"
	leftRow1    = "    /"
	leftRow2    = "   / "
	leftRow3    = "  /  "
	leftRow4    = " /   "
	leftRow5    = "/    "
	rightRow1   = "\\    "
	rightRow2   = " \\   "
	rightRow3   = "  \\  "
	rightRow4   = "   \\ "
	rightRow5   = "    \\"
	rightRow6   = "     \\"
	rightRow7   = "      \\"
	underbar    = "_____"
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
)

// indentOptions tracks the spacings used at a given depth and tree height for a given node width.
type indentOptions struct {
	// prefixPadding is how much spacing to start the beginning of a line with.
	// It is 2*height_from_bottom + 1 (where the height is 0 based).
	prefixPadding int

	// intraNodePadding is the spacing between the left and right legs of the tree.
	intraNodePadding int

	// interTreePadding is the spacing between each set of trees at this level.
	interTreePadding int

	// shoulderPadding is how much lateral filler we need between the top of a leg
	// and the levels node. (rather than growing the diagonal 2^n vertically, we
	// limit it to 5 vertical levels and then go sideways to make up the space.)
	shoulderPadding int

	// legDepth tracks how many of the vertical levels of the legs to show.
	legDepth int
}

var (
	// treeIndents is a mapping of distance from bottom level being rendered
	// to that level's indent options. By virtue of how wide things start to
	// get at 5 levels of tree, we don't plan to support more than that as text.
	treeIndents = map[int]indentOptions{
		0: indentOptions{
			prefixPadding:    0,
			intraNodePadding: 1,
			interTreePadding: 1,
			shoulderPadding:  0,
			legDepth:         4,
		},
		1: indentOptions{
			prefixPadding:    0,
			intraNodePadding: 1,
			interTreePadding: 1,
			shoulderPadding:  0,
			legDepth:         5,
		},
		2: indentOptions{
			prefixPadding:    2,
			intraNodePadding: 1,
			interTreePadding: 5,
			shoulderPadding:  0,
			legDepth:         5,
		},
		3: indentOptions{
			prefixPadding:    4,
			intraNodePadding: 1,
			interTreePadding: 9,
			shoulderPadding:  2,
			legDepth:         5,
		},
		// For a tree of height 5, this is the root node, so there is no
		// more tree or node padding above it so those values are 0.
		4: indentOptions{
			prefixPadding:    8,
			intraNodePadding: 1,
			interTreePadding: 0,
			shoulderPadding:  6,
			legDepth:         5,
		},
	}
)

// RenderMode is an enum for potential outputs when dumping or rendering the trees.
type RenderMode int

// Set of current render modes.
const (
	ModeASCII RenderMode = iota
	ModeSVG
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

	height := t.Height()
	node := t
	nodes := []BinaryTree[T]{node}
	var nextNodes []BinaryTree[T]
	depthFrom := height - 1
	indentOpts := treeIndents[depthFrom]

	// First pass starts with the root node, then we go into the loop of
	// legs and nodes until we are all done.
	outputNodes(nodes, indentOpts, &buf, depthFrom)

	for depthFrom > 0 {
		nextNodes = generateLevelsNodes(nodes)
		outputLegs(nextNodes, indentOpts, &buf, depthFrom)

		depthFrom--
		indentOpts = treeIndents[depthFrom]
		nodes = nextNodes
		outputNodes(nodes, indentOpts, &buf, depthFrom)
	}

	return buf.String()
}

// writeLeg replaces the boilerplate with a simple helper.
func writeLeg[T constraints.Ordered](leg BinaryTree[T], legString string, buf *bytes.Buffer) {
	if leg != nil {
		buf.WriteString(legString)
	} else {
		buf.WriteString(indent)
	}
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

// outputLegs does the boring bits of printing out visible or missing legs and the
// appropriate spacings between each one.
//
// TODO(rsned): Not visible in the end, but would be nice to not print the
// final appended spacing on the last potential tree leg.
// TODO(rsned): Minor enhancement would be to be able to skip the rest of a line
// when we get to the last actual leg piece of a row and the remaining parts
// would all just be spacings.
func outputLegs[T constraints.Ordered](nodes []BinaryTree[T], opts indentOptions, buf *bytes.Buffer, depthFrom int) {

	lastNode := lastNonNilNode(nodes)
	for i, ll := range leftLegs[:opts.legDepth] {
		buf.WriteString(indentFull[:elementWidth*opts.prefixPadding])
		for j := 0; j < len(nodes); j++ {
			if j > lastNode {
				break
			}
			writeLeg(nodes[j], ll, buf)

			// If this level has lateral legs, put in blanks to cover.
			buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])

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
			buf.WriteString(indentFull[:elementWidth*opts.intraNodePadding])

			// If this level has lateral leg elements, put in blanks to cover.
			buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])

			writeLeg(nodes[j], rightLegs[i], buf)

			// For all but the final node in the list.
			if j != len(nodes)-1 {
				// Spacing between subtrees.
				buf.WriteString(indentFull[:elementWidth*opts.interTreePadding])
			}
		}
		buf.WriteString("\n")
	}
}

// outputNodes writes out all the nodes and metadata at this level.
func outputNodes[T constraints.Ordered](nodes []BinaryTree[T], opts indentOptions, buf *bytes.Buffer, depthFrom int) {
	lastNode := lastNonNilNode(nodes)

	// Nodes.
	buf.WriteString(indentFull[:elementWidth*opts.prefixPadding])
	for j, n := range nodes {
		// This indent lines up with the space used to print left leg lines.
		// At the lowest level, there are no legs, so this leg padding needs to go.
		if depthFrom != 0 {
			buf.WriteString(indent)
		}

		// Higher up levels have lines that go sideways to keep the tree
		// reasonably sized.
		if n != nil && n.HasLeft() {
			buf.WriteString(underbarFull[:elementWidth*opts.shoulderPadding])
		} else {
			buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])
		}

		// The actual node value.
		if n != nil {
			buf.WriteString(fmt.Sprintf(nodeFmtT, n.Value()))
		} else {
			buf.WriteString(indent)
		}

		if n != nil && n.HasRight() {
			buf.WriteString(underbarFull[:elementWidth*opts.shoulderPadding])
		} else {
			buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])
		}
		// If this is the last node, skip all the remaining trailing padding.
		if j >= lastNode {
			break
		}

		// This indent lines up with the right leg lines.
		if depthFrom != 0 {
			buf.WriteString(indent)
		}
		buf.WriteString(indentFull[:elementWidth*opts.interTreePadding])
	}
	buf.WriteString("\n")

	if !levelHasMetadata(nodes) {
		return
	}

	// Add metadata print
	buf.WriteString(indentFull[:elementWidth*opts.prefixPadding])
	for j, n := range nodes {
		// This indent lines up with the left leg lines.
		if depthFrom != 0 {
			buf.WriteString(indent)
		}
		buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])
		if n != nil {
			buf.WriteString(fmt.Sprintf(nodeMetaFmt, n.Metadata()))
		} else {
			buf.WriteString(indent)
		}
		// If this is the last node, skip all the remaining trailing padding.
		if j >= lastNode {
			break
		}

		buf.WriteString(indentFull[:elementWidth*opts.shoulderPadding])
		// This indent lines up with the right leg lines.
		if depthFrom != 0 {
			buf.WriteString(indent)
		}
		buf.WriteString(indentFull[:elementWidth*opts.interTreePadding])
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

// centerString centers the given string into the target size adjusting the whitespace at
// either end as needed.
//
// This method assumes an output width of 50 or less for the purpose of this file.
func centerString(s string, width int) string {
	s = strings.TrimSpace(s)
	l := len(s)

	// For now, there is no attempt to truncate or elide longer values.
	if l >= width {
		return s
	}

	lPad := (width - l) / 2
	rPad := width - l - lPad

	const spaces = "                                                                     "
	return fmt.Sprintf("%s%s%s", spaces[0:lPad], s, spaces[0:rPad])
}
