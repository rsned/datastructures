package tree

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

var (
	treeSizeUpperLimit = flag.Int("tree_size_upper_limit", 500000,
		"What is the upper bound on how many values to insert when benchmarking. "+
			"(Generally this would be used when you want to make sure new code "+
			"is working by testing up to a smaller upper limit without waiting "+
			"for a full run of benchmarks.)")

	treeTypeFilter = flag.String("tree_type_filter", "",
		"Flag to restrict benchmark runs to one type of tree. String should match "+
			"the type name case-insensitively. e.g., avl or AvL or AVL would all "+
			"be matched to the AVL tree type.")
)

// To cut out some of the timing variability of benchmark functions, pre-create
// and sort a large set of random values.
var (
	testIntVals              []int
	testIntValsSorted        []int
	testIntValsReverseSorted []int
)

const (
	// limit is the maximum amount of random values to pregenerate for
	// the benchmarks to run against. This value also limits the value
	// in the flag "tree_size_upper_limit". If the flag is larger, then
	// it is limited to this amount.
	limit = 2000000
)

func init() {
	testIntVals = make([]int, limit)
	testIntValsSorted = make([]int, limit)
	testIntValsReverseSorted = make([]int, limit)

	for i := 0; i < limit; i++ {
		testIntVals[i] = rand.Int()
	}

	copy(testIntValsSorted, testIntVals)
	copy(testIntValsReverseSorted, testIntVals)

	sort.Ints(testIntValsSorted)
	sort.Sort(sort.Reverse(sort.IntSlice(testIntValsReverseSorted)))
}

var (
	// Define a series of increasing size amounts for benchmarking.
	insertSteps = []int{
		1000,
		5000,
		10000,
		15000,
		20000,
		25000,
		50000,
		100000,
		250000,
		500000,
	}
)

// newTreeFunc is used by the benchmarks to instantiate a new instance
// of some Tree type to be used in a benchmark run.
type newTreeFunc[T constraints.Ordered] func() Tree[T]

// newBSTTree creates a new BST.
func newBSTTree[T constraints.Ordered]() Tree[int] {
	return &BST[int]{}
}

// newAVLTree creates a new AVL.
func newAVLTree[T constraints.Ordered]() Tree[int] {
	return &AVL[int]{}
}

// To run the Tree Insert Benchmarks use this command with
// the desired number of run repetitions:
//
// go test . --test.benchmem --test.bench="BenchmarkTreeInsert" --count=n
//

// BenchmarkTreeInsert is a harness to benchmark the Insert method on all
// supported tree types.
func BenchmarkTreeInsert(b *testing.B) {
	examples := []struct {
		name   string
		sorted bool
		tree   newTreeFunc[int]
	}{
		{
			name: "BST",
			tree: newBSTTree[int],
		},
		{
			name: "AVL",
			tree: newAVLTree[int],
		},
	}

	for _, example := range examples {
		// Check if the user requested filtering on the benchmark.
		if *treeTypeFilter != "" &&
			strings.EqualFold(example.name, *treeTypeFilter) {
			continue
		}

		for _, n := range insertSteps {
			// Skip any tests that are outside the limit.
			if n > *treeSizeUpperLimit {
				break
			}
			vals := testIntVals[:n]

			b.Run(fmt.Sprintf("%s-%06d", example.name, n),
				func(b *testing.B) {
					tree := example.tree()
					for i := 0; i < b.N; i++ {
						tree.Insert(vals[i%n])
					}
				})
		}
	}
}

// BenchmarkTreeSearch is a harness to benchmark the Search method
// on all supported tree types.
func BenchmarkTreeSearch(b *testing.B) {
	examples := []struct {
		name string
		tree newTreeFunc[int]
	}{
		{
			name: "BST",
			tree: newBSTTree[int],
		},
		{
			name: "AVL",
			tree: newAVLTree[int],
		},
	}

	for _, example := range examples {
		// Check if the user requested filtering on the benchmark.
		if *treeTypeFilter != "" &&
			strings.EqualFold(example.name, *treeTypeFilter) {
			continue
		}

		for _, n := range insertSteps {
			b.StopTimer()
			// Skip any tests that are outside the limit.
			if n > *treeSizeUpperLimit {
				break
			}
			vals := testIntVals[:n]
			// Build and fill the tree before starting the benchmark.
			tree := example.tree()
			for i := 0; i < n; i++ {
				tree.Insert(vals[i])
			}

			b.StartTimer()
			b.Run(fmt.Sprintf("%s-%06d", example.name, n),
				func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_ = tree.Search(vals[i%n])
					}
				})
		}
	}
}
