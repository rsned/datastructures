package tree

import "golang.org/x/exp/constraints"

// Options contains the various settings used in these tree functions.
type Options struct {
	// ignoreDuplicates indicates of duplicate values
	ignoreDuplicates bool

	// fpTolerance is used to set floating point tolerance for equality
	// comparisons.
	fpTolerance float64
}

func defaultOptions() *Options {
	return &Options{
		ignoreDuplicates: true,
		fpTolerance:      1e-15,
	}
}

// treeOptionFunc is a function to set options for use in variadic opt params.
type OptionFunc func(c *Options)

// IgnoreDuplicates tells the tree function that a duplicate value in an
// operation should be ignored. (Such as when joining two Trees)
func IgnoreDuplicates(ignore bool) OptionFunc {
	return func(o *Options) {
		o.ignoreDuplicates = ignore
	}
}

// FloatingPointTolerance sets the tolerance when compariong Floating Point
// values in tree operations.
func FloatingPointTolerance(tol float64) OptionFunc {
	return func(o *Options) {
		o.fpTolerance = tol
	}
}

// Join attempts to merge the given trees following the given options (if any).
// If the joining was unsuccessful for any reason, the resulting Tree should
// not be used, and false will be returned.
//
// Options can include things like what strategy to use when encountering
// duplicate values, hints or requirements on type of output tree, etc.
func Join[T constraints.Ordered](a, _ Tree[T], opts ...OptionFunc) (Tree[T], bool) {
	treeOpts := defaultOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	// TODO(rsned): Implement this.
	return a, false
}

// Split splits the Tree into two trees such that first tree returned constains
// the values up to and including the split point, and the second tree the
// remainder. The output Trees will be of the same underlying type as the input.
//
// If the value falls between two nodes in the tree, then tree one will end at
// the value closest without exceeding the given value.
//
// The resulting Trees are NOT guaranteed to be balanced or optimal.
//
// Your right to be foolish is supported. For example splitting on a value
// outside the trees limits will give back the original tree and a nil tree.
func Split[T constraints.Ordered](t Tree[T], _ T) (Tree[T], Tree[T]) {
	// TODO(rsned): Implement this.
	return t, t
}

// Prune removes the whole subtree that is homed at val.
//
// Your right to be foolish is supported. For example pruning on the root
// value will give back an empty tree.
func Prune[T constraints.Ordered](t Tree[T], _ T) Tree[T] {
	// TODO(rsned): Implement this.
	return t
}

// Rebalance attempts to perform some after-market rebalancing on a tree.
//
// For types that normally include some type of balancing, they may
// short-circuit this call.
func Rebalance[T constraints.Ordered](t Tree[T], _ T) Tree[T] {
	// TODO(rsned): Implement this.
	return t
}

// Convert attempts to convert the given tree into a tree of a type specified
// in the options. For example, Convert a BST to an AVL tree.
//
// More discourse will follow on how different combinations of options will
// be handled such as specifying multiple tree types, etc.
//
//	opts: Underlying type
func Convert[T constraints.Ordered](t Tree[T], opts ...OptionFunc) Tree[T] {
	treeOpts := defaultOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	// TODO(rsned): Implement this.
	return t
}

// ToSlice converts the tree to a slice in natural order.
func ToSlice[T constraints.Ordered](_ Tree[T]) []T {
	var x []T

	// TODO(rsned): Implement this.
	return x
}

// Equal reports if the two trees containt the same nodes in the same structure.
//
// It is possible for two different types of Binary Trees, for example, to have
// the same nodes and the same structure.
//
// e.g., Tree A, a BST:
//
//	  5
//	 / \
//	2   8
//
// and Tree B, an AVL tree:
//
//	  5
//	 / \
//	2   8
//
// Are equal because the trees have the same nodes and structure.
//
// Whereas if Tree A was:
//
//	    8
//	   /
//	  5
//	 /
//	2
//
// It would be equivalent, but not equal, because it has the same node
// values in an In-Order traversal, but the structure is different.
//
// This function supports changing the tolerance for floating point comparisons.
func Equal[T constraints.Ordered](a, b Tree[T], opts ...OptionFunc) bool {
	treeOpts := defaultOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	// TODO(rsned): Once other types of Trees exist besides BinaryTree,
	// enhance this to choose the appropriate equality.
	aBinary, aOk := a.(BinaryTree[T])
	bBinary, bOk := b.(BinaryTree[T])
	if !aOk || !bOk {
		return false
	}

	return BinaryTreesEqual(aBinary, bBinary)
}

// Equivalent reports if the two trees have the same node values in the same
// order. This is essentially reporting if the two trees have the same In-Order
// traversal outputs, but not caring about the underlying structure or
// implementation.
//
// See the description for Equal for some examples of this.
//
// e.g., Tree A, a BST:
//
//		    8
//	       / \
//		  5  15
//	     /  /  \
//	    2  13  17
//
// and Tree B, a B-Tree
//
//		      13
//	         /  \
//		2|5|8    15|17
//
// are equivalent because they have the same node values in the same order.
//
// This function supports changing the tolerance for floating point comparisons.
func Equivalent[T constraints.Ordered](a, b Tree[T], opts ...OptionFunc) bool {
	treeOpts := defaultOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	// TODO(rsned): Once other types of Trees exist besides BinaryTree,
	// enhance this to choose the appropriate equality.
	aBinary, aOk := a.(BinaryTree[T])
	bBinary, bOk := b.(BinaryTree[T])
	if !aOk || !bOk {
		return false
	}

	return BinaryTreesEquivalent(aBinary, bBinary)
}

// Summarize takes a tree and returns comprehensive statistics and analysis
// about the tree structure, balance, performance characteristics, and more.
func Summarize[T constraints.Ordered](_ Tree[T]) *Summary[T] {
	var summary *Summary[T]

	// TODO(rsned): Implement this.

	return summary
}
