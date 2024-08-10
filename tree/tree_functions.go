package tree

import "golang.org/x/exp/constraints"

// TreeOptions contains the various settings used in these tree functions.
type TreeOptions struct {
	// ignoreDuplicates indicates of duplicate values
	ignoreDuplicates bool

	// FPTolerance is used to set floating point tolerance for equality
	// comparisons.
	fpTolerance float64
}

func defaultTreeOptions() *TreeOptions {
	return &TreeOptions{
		ignoreDuplicates: true,
		fpTolerance:      1e-15,
	}
}

// treeOptionFunc is a function to set options for use in variadic opt params.
type treeOptionFunc func(c *TreeOptions)

// IgnoreDuplicates tells the tree function that a duplicate value in an
// operation should be ignored. (Such as when joining two Trees)
func IgnoreDuplicates(ignore bool) treeOptionFunc {
	return func(o *TreeOptions) {
		o.ignoreDuplicates = ignore
	}
}

// FloatingPointTolerance sets the tolerance when compariong Floating Point
// values in tree operations.
func FloatingPointTolerance(tol float64) treeOptionFunc {
	return func(o *TreeOptions) {
		o.fpTolerance = tol
	}
}

// Clone returns a complete new copy of the given tree.
func Clone[T constraints.Ordered](t Tree[T]) Tree[T] {
	return t
}

// Join combines the given trees using the options (if any).
//
// Options can include things like what strategy to use when encountering
// duplicate values, hints or reqeuirements on type of output tree, etc.
func Join[T constraints.Ordered](a, b Tree[T], opts ...treeOptionFunc) Tree[T] {
	treeOpts := defaultTreeOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}
	return a
}

// Split splits the Tree into two trees such that first tree returned constains
// the values up to and including the split point, and the second tree the
// remainder. The output Trees will be of the same underlying type as the input.
//
// If the value falls between two nodes in the tree, then tree one will end at
// the value closest without exceeding the given value.
func Split[T constraints.Ordered](t Tree[T], val T) (Tree[T], Tree[T]) {
	return t, t
}

// Prune removes the whole subtree that is homed at val.
func Prune[T constraints.Ordered](t Tree[T], val T) Tree[T] {
	return t
}

// Rebalance attempts to perform some rebalancing on a tree.
//
// Not all types need it so those types may short-circuit this.
func Rebalance[T constraints.Ordered](t Tree[T], val T) Tree[T] {
	return t
}

// Convert attempts to convert the given tree into a tree of a type specified
// in the options.
//
// More discourse will follow on how different combinations of options will
// be handled such as specifying multiple tree types, etc.
//
//	opts: Underlying type
func Convert[T constraints.Ordered](t Tree[T], opts ...treeOptionFunc) Tree[T] {
	treeOpts := defaultTreeOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	return t
}

// ToSlice converts the tree to a slice in natural order.
func ToSlice[T constraints.Ordered](t Tree[T]) []T {
	var x []T
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
// values in an In Order traversal, but the structure is different.
//
// This function supports changing the tolerance for floating point comparisons.
func Equal[T constraints.Ordered](a, b Tree[T], opts ...treeOptionFunc) bool {
	treeOpts := defaultTreeOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	return false
}

// Equivalent reports if the two trees have the same node values in the same order.
// This is essentially reporting if the two trees have the same In Order traversal
// outputs, but not caring about the underlying structure.
//
// See the description for Equal for examples of this.
//
// This function supports changing the tolerance for floating point comparisons.
func Equivalent[T constraints.Ordered](a, b Tree[T], opts ...treeOptionFunc) bool {
	treeOpts := defaultTreeOptions()
	for _, opt := range opts {
		opt(treeOpts)
	}

	return false
}
