package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// PrintTreeASCII outputs the Tree up to a reasonable number of levels deep
// for the purpose of aiding in testing and debugging.
//
// This method takes a Tree[T] with an explicit assumption the values in the
// tree are under a predefined width when formatted and printed. e.g. -21, 123,
// 7, 6.02e23, "Yellow", etc. Values wider than the cutoff may work fine, but
// no testing or quality checks have been done.
//
// An optional label is output before the tree contents with a blank line
// separator between the label and the tree.
//
// TODO(rsned): Add an elideLevel int param to do the eliding dynamically.
func PrintTreeASCII[T constraints.Ordered](label string, t Tree[T]) string {
	var buf bytes.Buffer

	if len(label) > 0 {
		buf.WriteString(label + "\n\n")
	}

	if isTreeNil(t) {
		buf.WriteString("nil tree\n")

		return buf.String()
	}

	// Specific tree types need more specific type casting to print.

	switch typ := t.(type) {
	case *BST[T]:
		buf.WriteString(PrintBinaryTreeASCII("", typ.root))
	case *AVL[T]:
		buf.WriteString(PrintBinaryTreeASCII("", typ.root))
	case *RedBlack[T]:
		buf.WriteString(PrintBinaryTreeASCII("", typ.root))
	default:
		buf.WriteString(fmt.Sprintf("unknown tree type: %q\n", typ))
	}

	return buf.String()
}
