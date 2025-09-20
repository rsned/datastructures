package tree

import (
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
	indent      = "     "
	nodeFmt     = "%-3v"
	nodeMetaFmt = "%5s"

	// These are the constant strings used for rendering leg segements of
	// the given depth.
	leftRow1  = "/"
	leftRow2  = "/ "
	leftRow3  = "/  "
	leftRow4  = "/   "
	leftRow5  = "/    "
	leftRow6  = "/     "
	leftRow7  = "/      "
	rightRow1 = "\\"
	rightRow2 = " \\"
	rightRow3 = "  \\"
	rightRow4 = "   \\"
	rightRow5 = "    \\"
	rightRow6 = "     \\"
	rightRow7 = "      \\"
	underbar  = "_____"
)

// maxPadding is how long to make the pad strings we substring against.
const maxPadding = 512

var (
	// leftLegs is a slice of the angled leg strings in order to be
	// iterated over at each level to make the render function cleaner.
	leftLegs = []string{
		leftRow1,
		leftRow2,
		leftRow3,
		leftRow4,
		leftRow5,
		leftRow6,
		leftRow7,
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
	underbarFull = strings.Repeat(underbar, maxPadding)
	indentFull   = strings.Repeat(indent, maxPadding)

	// The remaining are defined separately even though they use the same
	// string normally. This allows debugging by tweaking the given character
	// so we can see it in outputs.
	prefixPad   = strings.Repeat(" ", maxPadding) // Switch to 'P' for debug.
	shoulderPad = strings.Repeat(" ", maxPadding) // Switch to 'S' for debug.
	interPad    = strings.Repeat(" ", maxPadding) // Switch to 'I' for debug.
	intraPad    = strings.Repeat(" ", maxPadding) // Switch to 'i' for debug.
	otherPad    = strings.Repeat(" ", maxPadding) // Switch to '#' for debug.
	otherPad2   = strings.Repeat(" ", maxPadding) // Switch to '$' for debug.
	legPad      = strings.Repeat(" ", maxPadding) // Switch to 'L' for debug.
)

// indentOptions tracks the spacings used at a given depth and tree height for
// a given node width.
//
// Using the following tree as an example to point out which fields represent
// which parts of this struct.
//
//	P = prefixPad
//	S = shoulderPad
//	I = interPad
//	i = intraPad AKA node width
//	L = legPad (Usually only 1, but some variations of the tool used more
//	    vertical space for higher levels to shorten shoulders)
//
// | <- Edge of output area.
// |
// |                 __________500__________
// |                /                       \
// |          ___250___                   ___750___
// |         /         \                 /         \
// |      125           375           625           875
// |     /   \         /   \         /   \         /   \
// |  100     187   225     425   123     321   123     999
//
// | <- Edge of output area.
// |
// |PPPPPPPPPPPPPPPPL__________500__________L
// |                /SSSSSSSSSSiiiSSSSSSSSSS\
// |PPPPPPPPPL___250___                   ___750___
// |PPPPPPPPP/SSSiiiSSS\                 /SSSiiiSSS\
// |PPPPPL125LSSSiiiSSSL375LIIIIIIIIIL625LSSSiiiSSSL875
// |PPPPP/iii\IIIIIIIII/iii\         /   \         /   \
// |PP100LiiiL187III225LiiiL425III123LiiiL321III123LiiiL999
// |
type indentOptions struct {
	// nodeWidth is how wide in number of spaces this node value will be.
	nodeWidth int

	// prefixPadding is how much spacing to start the beginning of a line with.
	// This is to shift the entire output over. It is not the leading space on
	// rows higher up the tree that do not start flush with the left side of the
	// output area.
	//
	// Element 'P' in the diagram above.
	//
	// This is measured in spaces.
	prefixPadding int

	// intraNodePadding is the spacing between the left and right legs of the
	// tree. Generally this matches the nodeWidth, but it can be less if the
	// tree is trying to be more densely rendered. (Or more to be more
	// cushion-y which will bump up the shoulderPadding at each level.)
	//
	// Element 'i' in the diagram above.
	//
	// This is measured in spaces.
	intraNodePadding int

	// interTreePadding is the spacing between each set of trees at this level.
	// This is usually less than the intraNodePadding and nodeWidth and
	// primarily affects the shoulder spacing on higher up levels.
	//
	// Element 'I' in the diagram above.
	//
	// This is measured in units of spaces.
	interTreePadding int

	// shoulderPadding is how much lateral filler we need between the top of a
	// leg and the current level's node text. (rather than growing the diagonal
	// 2^n vertically, we limit it to one vertical level and then go sideways
	// to make up the space needed to get it into position.)
	//
	// Element 'S' in the diagram above.
	//
	// This is in spaces.
	shoulderPadding int

	// legDepth tracks how tall the legs are between the current level and the
	// level above it. Mostly this value is 1 because we are focusing on
	// compact vertical spacing. Larger values can generate a more artisinal
	// feel.
	//
	// There are no hard upper limits, but going beyond 5 really strains the
	// visual sensibilities.
	legDepth int
}

// RenderMode is an enum for output formats when printing out trees.
type RenderMode int

// Set of current render modes.
const (
	ModeASCII RenderMode = iota // Also the default value.
	ModeSVG

	// TODO(rsned): Add more modes?
)

// RenderBinaryTree returns the given tree in the given mode rendered into string form.
func RenderBinaryTree[T constraints.Ordered](t BinaryTree[T], _ int, mode RenderMode) string {
	switch mode {
	case ModeASCII:
		return PrintBinaryTreeASCII("", t)
	case ModeSVG:
		return "SVG method not implemented yet"
	default:
		return "Method not implemented yet"
	}
}

func indentOptsForNodeWidth(width int) indentOptionsMap {
	// TODO(rsned): Add check for maxSupportedWidth to prevent crashes.
	return binaryTreeSpacingData[width-(width+1)%2]
}

// generateLevelsNodes returns a potentially sparse slice of Nodes at the
// next level in the tree based on the current slice of tree Nodes. Nil
// Nodes and any nil children are replaced with nils as placeholders in
// the output.
//
// This slice should be a sparse slice of size 2*N, where N is the number
// of nodes from the previous level.
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
			// If the current node was nil (i.e. an unbalanced
			// tree in progress), we output two blank filler entries
			// for its non-existent children.
			nodes = append(nodes, nil)
			nodes = append(nodes, nil)
		}
	}

	return nodes
}

// lastNonNilNode walks backward through the list looking for the farthest
// right non-nil node to know when the current row can break early in
// processing.
func lastNonNilNode[T constraints.Ordered](nodes []BinaryTree[T]) int {
	for i := len(nodes) - 1; i >= 0; i-- {
		if nodes[i] != nil {
			return i
		}
	}

	return 0
}

// centerString centers the given string into the target size adjusting
// the space at either end as needed with the given pad character.
func centerString(s, padChar string, width int) string {
	s = strings.TrimSpace(s)
	l := len(s)

	// For now, there is no attempt to truncate or elide longer values.
	if l >= width {
		return s
	}

	// Calculate padding: for uneven splits, put the extra space on the right
	// This matches the test expectations where width 2 with "a" becomes "a "
	lPad := (width - l) / 2
	rPad := width - l - lPad

	// Create padding strings using the provided padChar
	leftPadding := strings.Repeat(padChar, lPad)
	rightPadding := strings.Repeat(padChar, rPad)

	return leftPadding + s + rightPadding
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
		widestValue: 0,
	}

	// things we want to find out:
	// max height
	// width of largest value
	// lopsidedness / skew   e.g. is this only a one sided binary tree?

	// Walk the tree getting all values and printing them as strings in order
	// to find the widest, min, max, mean, etc.
	ch := tree.Traverse(TraverseInOrder)

	var widest int
	var minVal, maxVal T

	for {
		val, ok := <-ch
		if ok {
			s := fmt.Sprintf("%v", val)
			if len(s) > widest {
				widest = len(s)
			}
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		} else {
			break
		}
	}

	stats.widestValue = widest

	return stats
}
