package tree

import (
	"math/rand"
	"slices"
	"sort"
	"testing"
)

// randomBSTValLimit sets an arbitrary limit on the value of the nodes
// to keep trees readable when printed out.
const randomBSTValLimit = 999

// generateRandomBSTVals generates a slice of unique random values in the range
// (-limit, limit) that can be used to populate a binary search tree.
func generateRandomBSTVals(t *testing.T, seed int64, numNodes int, limit int) []int {
	t.Helper()

	r := rand.New(rand.NewSource(seed))
	vals := make([]int, 0, numNodes)
	seen := make(map[int]bool, numNodes)

	for len(vals) < numNodes {
		// Generate a random value in the range
		// (-limit, +limit)
		val := r.Intn(2*limit) - limit

		// Only add if we haven't seen this value before
		if !seen[val] {
			seen[val] = true
			vals = append(vals, val)
		}
	}

	return vals
}

// generateLevelOrderTraverseValues takes a sorted slice of values and a
// splitPoint to start as the root node and creates a slice of the values as
// they would be traversed in level-order so that they can be fed into a
// tree in the correct order to obtain the desired form.
func generateLevelOrderTraverseValues(t *testing.T, vals []int, midPoint int) []int {
	t.Helper()

	numVals := len(vals)
	if numVals == 0 {
		return []int{}
	}

	if numVals == 1 {
		return []int{vals[0]}
	}

	levelOrderVals := make([]int, 0, numVals)

	// rangeInfo represents a range in the vals slice with its midpoint
	type rangeInfo struct {
		start int // inclusive
		end   int // exclusive
		mid   int // the midpoint index within this range
	}

	// Use BFS to traverse level by level
	queue := []rangeInfo{{start: 0, end: numVals, mid: midPoint}}

	for len(queue) > 0 {
		// Process all nodes at current level
		levelSize := len(queue)
		for range levelSize {
			curr := queue[0]
			queue = queue[1:]

			// Add the midpoint value to result
			levelOrderVals = append(levelOrderVals, vals[curr.mid])

			// Calculate left child range [start, mid)
			if curr.mid > curr.start {
				leftStart := curr.start
				leftEnd := curr.mid
				leftMid := (leftStart + leftEnd) / 2
				queue = append(queue, rangeInfo{start: leftStart, end: leftEnd, mid: leftMid})
			}

			// Calculate right child range (mid, end]
			if curr.mid+1 < curr.end {
				rightStart := curr.mid + 1
				rightEnd := curr.end
				// Midpoint is calculated as (curr.mid + rightEnd) / 2
				// This gives a more balanced distribution
				rightMid := (curr.mid + rightEnd) / 2
				queue = append(queue, rangeInfo{start: rightStart, end: rightEnd, mid: rightMid})
			}
		}
	}

	return levelOrderVals
}

// generateStructuredBinaryTree is a helper that creates a BST with a specific structure
// by generating random values, sorting them, optionally applying level-order transformation,
// and optionally reversing before insertion.
func generateStructuredBinaryTree(t *testing.T, seed int64, numNodes int, useLevelOrder bool, midpointFunc func(int) int, reverse bool) Tree[int] {
	t.Helper()

	vals := generateRandomBSTVals(t, seed, numNodes, randomBSTValLimit)
	sort.Ints(vals)

	if useLevelOrder {
		midPoint := midpointFunc(len(vals))
		vals = generateLevelOrderTraverseValues(t, vals, midPoint)
	} else if reverse {
		slices.Reverse(vals)
	}

	tree := NewBST[int]()
	for _, val := range vals {
		tree.Insert(val)
	}

	return tree
}

// generateRandomBinaryTree generates a BST with the given number of nodes with
// random values and NO attempt at ordering or balance or evenness.
func generateRandomBinaryTree(t *testing.T, seed int64, numNodes int) Tree[int] {
	t.Helper()

	tree := NewBST[int]()
	for _, v := range generateRandomBSTVals(t, seed, numNodes, randomBSTValLimit) {
		tree.Insert(v)
	}

	return tree
}

// generateBalancedBinaryTree generates a BST that attempts to be
// balanced or reasonably close to balanced.
func generateBalancedBinaryTree(t *testing.T, seed int64, numNodes int) Tree[int] {
	t.Helper()

	return generateStructuredBinaryTree(t, seed, numNodes, true, func(n int) int { return n / 2 }, false)
}

// generatedSkewedBinaryTree generates a BST that attempts to be
// skewed to either the left or right.
func generatedSkewedBinaryTree(t *testing.T, seed int64, numNodes int, left bool) Tree[int] {
	t.Helper()

	// If we are going to skew to the left, split the values at the ~3/4 point.
	// Otherwise if we are going to the right, split the values at ~1/4 point.
	return generateStructuredBinaryTree(t, seed, numNodes, true, func(n int) int {
		if left {
			return 3 * n / 4
		}

		return n / 4
	}, false)
}

// generateDegenerateBinaryTree generates a BST that attempts to be
// degenerate or nearly degenerate (usually a one-legged linear tree).
func generateDegenerateBinaryTree(t *testing.T, seed int64, numNodes int, left bool) Tree[int] {
	t.Helper()

	return generateStructuredBinaryTree(t, seed, numNodes, false, nil, left)
}

// TestGenerateRandomBSTVals verifies that generateRandomBSTVals produces
// unique values within the expected range.
func TestGenerateRandomBSTVals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		seed     int64
		numNodes int
		limit    int
	}{
		{
			name:     "small tree",
			seed:     42,
			numNodes: 10,
			limit:    20,
		},
		{
			name:     "medium tree",
			seed:     123,
			numNodes: 50,
			limit:    50,
		},
		{
			name:     "large tree",
			seed:     999,
			numNodes: 200,
			limit:    200,
		},
		{
			name:     "single node",
			seed:     1,
			numNodes: 1,
			limit:    1,
		},
		{
			name:     "edge case - near limit",
			seed:     777,
			numNodes: 500,
			limit:    randomBSTValLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vals := generateRandomBSTVals(t, tt.seed, tt.numNodes, tt.limit)

			// Check we got the right number of values
			if len(vals) != tt.numNodes {
				t.Errorf("generateRandomBSTVals() returned %d values, want %d", len(vals), tt.numNodes)
			}

			// Check for duplicates using a map
			seen := make(map[int]bool)
			for _, v := range vals {
				if seen[v] {
					t.Errorf("generateRandomBSTVals() produced duplicate value: %d", v)

					break
				}
				seen[v] = true
			}

			// Check all values are within valid range
			for _, v := range vals {
				if v < -tt.limit || v >= tt.limit {
					t.Errorf("generateRandomBSTVals() produced value %d outside range [%d, %d)",
						v, -tt.limit, tt.limit)

					break
				}
			}
		})
	}
}

// TestGenerateLevelOrderTraverseValues tests the generateLevelOrderTraverseValues function
// with various input sizes and midpoint positions.
func TestGenerateLevelOrderTraverseValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		vals     []int
		midPoint int
		want     []int
	}{
		{
			name:     "empty slice",
			vals:     []int{},
			midPoint: 0,
			want:     []int{},
		},
		{
			name:     "single value",
			vals:     []int{5},
			midPoint: 0,
			want:     []int{5},
		},
		{
			name:     "two values - midpoint at 0",
			vals:     []int{1, 2},
			midPoint: 0,
			want:     []int{1, 2},
		},
		{
			name:     "two values - midpoint at 1",
			vals:     []int{1, 2},
			midPoint: 1,
			want:     []int{2, 1},
		},
		{
			name:     "three values - midpoint at 1 (middle)",
			vals:     []int{1, 2, 3},
			midPoint: 1,
			want:     []int{2, 1, 3},
		},
		{
			name:     "six values - midpoint at 3 (middle)",
			vals:     []int{1, 2, 3, 4, 5, 6},
			midPoint: 3,
			want:     []int{4, 2, 5, 1, 3, 6},
		},
		{
			name:     "ten values - midpoint at 5 (middle)",
			vals:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			midPoint: 5,
			want:     []int{6, 3, 8, 2, 4, 7, 9, 1, 5, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := generateLevelOrderTraverseValues(t, tt.vals, tt.midPoint)
			if len(got) != len(tt.want) {
				t.Errorf("generateLevelOrderTraverseValues() length = %v, want %v", len(got), len(tt.want))

				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("generateLevelOrderTraverseValues()[%d] = %v, want %v\nGot:  %v\nWant: %v",
						i, got[i], tt.want[i], got, tt.want)

					break
				}
			}
		})
	}
}

// TestGenerateLevelOrderTraverseValuesVariedMidpoints tests the function
// with various midpoint positions on a larger slice.
func TestGenerateLevelOrderTraverseValuesVariedMidpoints(t *testing.T) {
	t.Parallel()

	// Create a slice with at least 12 values for testing various midpoints
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	numVals := len(vals)

	tests := []struct {
		name        string
		midPoint    int
		description string
		want        []int
	}{
		{
			name:        "midpoint at 0% (index 0)",
			midPoint:    0,
			description: "completely right-skewed tree",
			want:        []int{1, 8, 5, 12, 3, 6, 10, 14, 2, 4, 7, 9, 11, 13, 15},
		},
		{
			name:        "midpoint at 5% (index ~0.75)",
			midPoint:    int(float64(numVals) * 0.05),
			description: "heavily right-skewed tree",
			want:        []int{1, 8, 5, 12, 3, 6, 10, 14, 2, 4, 7, 9, 11, 13, 15},
		},
		{
			name:        "midpoint at 25% (index ~3.75)",
			midPoint:    int(float64(numVals) * 0.25),
			description: "moderately right-skewed tree",
			want:        []int{4, 2, 10, 1, 3, 7, 13, 6, 8, 12, 14, 5, 9, 11, 15},
		},
		{
			name:        "midpoint at 75% (index ~11.25)",
			midPoint:    int(float64(numVals) * 0.75),
			description: "moderately left-skewed tree",
			want:        []int{12, 6, 14, 3, 9, 13, 15, 2, 4, 8, 10, 1, 5, 7, 11},
		},
		{
			name:        "midpoint at 95% (index ~14.25)",
			midPoint:    int(float64(numVals) * 0.95),
			description: "heavily left-skewed tree",
			want:        []int{15, 8, 4, 11, 2, 6, 10, 13, 1, 3, 5, 7, 9, 12, 14},
		},
		{
			name:        "midpoint at 100% (last index)",
			midPoint:    numVals - 1,
			description: "completely left-skewed tree",
			want:        []int{15, 8, 4, 11, 2, 6, 10, 13, 1, 3, 5, 7, 9, 12, 14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := generateLevelOrderTraverseValues(t, vals, tt.midPoint)
			if len(got) != len(tt.want) {
				t.Errorf("generateLevelOrderTraverseValues() length = %v, want %v", len(got), len(tt.want))

				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("%s: generateLevelOrderTraverseValues()[%d] = %v, want %v\nGot:  %v\nWant: %v",
						tt.description, i, got[i], tt.want[i], got, tt.want)

					break
				}
			}
		})
	}
}

// TODO(rsned): Remaining methods to test.
// traverseBinaryTree
