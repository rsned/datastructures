# Go Tree Data Structures

This package contains a variety of common ordered, rooted tree types implemented in Go using generics. It supports all types that satisfy `constraints.Ordered`, including integers, floating-point numbers, and strings.

There are implementations of a number of traditional Binary Trees (Binary Search Tree, AVL Tree, Red/Black tree, etc.) as well as future updates to include some N-ary (multiple values per node) trees like B-tree and B+-tree.

More complex or esoteric tree types are currently not in scope for this package.

The package contains comprehensive unit tests as well as a suite of Benchmarks that can be used to compare and contrast the various trees performances.

This implementation is not focused on ultimate performance, but rather on letting me work towards simplicity, readability, and ease of use.

## Currently Implemented Tree Types

### Binary Trees

#### Binary Search Tree (BST)


Simplest of all binary trees, a node, two children, no balancing. Simple insert/delete/search operations.


#### Adelson Velsky and Landis (AVL) Tree

A self-balancing binary tree where the height of any two child subtrees differs by no more than 1. Rebalancing happens automatically during insert and delete operations if the resulting tree is unbalanced.

#### Red-Black Tree
Approximately balanced binary tree with red/black node coloring that offers a good balance between insertion performance and search performance.


### N-ary Tree Types

In graph theory, an n-ary tree (for nonnegative integers n) is an ordered tree in which each node has no more than n children. 

#### B-Tree

A B-tree is a self-balancing tree data structure that maintains sorted data and allows searches, sequential access, insertions, and deletions in logarithmic time. 

#### B+-Tree

A B+ tree can be viewed as a B-tree in which each node contains only keys (not key–value pairs),
and to which an additional level is added at the bottom with linked leaves.


## Core Interfaces

### Tree Interface
All tree types implement the `Tree[T]` interface:

```go
type Tree[T constraints.Ordered] interface {
    Insert(v T) bool        // Insert a value
    Delete(v T) bool        // Delete a value
    Search(v T) bool        // Search for a value
    Height() int            // Get tree height
    Clone[T]() Tree[T]      // Create a deep copy
    Traverse(order TraverseOrder) <-chan T  // Traverse tree
}
```

### BinaryTree Interface
Binary tree nodes implement the `BinaryTree[T]` interface:

```go
type BinaryTree[T constraints.Ordered] interface {
    Tree[T]
    HasLeft() bool          // Check if left child exists
    HasRight() bool         // Check if right child exists
    Left() BinaryTree[T]    // Get left child
    Right() BinaryTree[T]   // Get right child
    Value() T               // Get node value
    Metadata() string       // Reports node-specific metadata
}
```

## Tree Utility Functions

The package provides some utility functions for tree operations:

- **`Equal[T](a, b Tree[T]) bool`**: Reports if two trees have the same values AND the same structure regardless of underlying type.  (e.g., A balanced AVL tree and a balanced BST would be considered Equal if they have the same values and structure.)
- **`Equivalent[T](a, b Tree[T]) bool`**: Reports if the two trees have the same values in the same order but ignoring tree structure.  (So a one-legged BST and and AVL could have be equivalent, but not Equal. Similarly, a B-Tree and a BST could be Equivalent but won't have the same structure.)

## Traversal Orders

The package supports multiple tree traversal methods:

```go
type TraverseOrder int
const (
    TraverseInOrder         // Left, Root, Right
    TraversePreOrder        // Root, Left, Right
    TraversePostOrder       // Left, Right, Root
    TraverseReverseOrder    // Right, Root, Left
    TraverseLevelOrder      // Breadth-first in order traversal
)
```

The `Traverse` method uses a `<-chan T` to return the values in the specified order. The channel is closed when the traversal is complete.

## Visualization

### ASCII Tree Display
Use `PrintBinaryTreeASCII` to visualize tree structure:

```go
tree := NewAVL[int]()
tree.Insert(10)
tree.Insert(5)
tree.Insert(15)
tree.Insert(12)

fmt.Println(PrintBinaryTreeASCII("My First AVL Tree", tree.Root()))
```

Output:
```
My First AVL Tree

     10
     (0)
     / \
    5   15
  (0)   (-1)
        /
       12
       (0)
```

## Usage Examples

### Basic Usage
```go
// Create an AVL tree for integers
tree := NewAVL[int]()

// Insert values
tree.Insert(10)
tree.Insert(5)
tree.Insert(15)
tree.Insert(3)
tree.Insert(7)

// Search for values
found := tree.Search(7)  // returns true
missing := tree.Search(12)  // returns false

// Delete values
deleted := tree.Delete(5)  // returns true

// Get tree height
height := tree.Height()  // returns height of tallest leg of the tree
```

### Tree Traversal
```go
// Traverse tree in different orders
for value := range tree.Traverse(TraverseInOrder) {
    fmt.Printf("%d ", value)
}
// Output: 3 7 10 15 

for value := range tree.Traverse(TraverseLevelOrder) {
    fmt.Printf("%d ", value)
}
// Output: 10 7 15 3
```

### Working with Different Types
```go
// String tree
stringTree := NewBST[string]()
stringTree.Insert("apple")
stringTree.Insert("banana")
stringTree.Insert("cherry")

// Float tree
floatTree := NewAVL[float64]()
floatTree.Insert(3.14)
floatTree.Insert(2.71)
floatTree.Insert(1.41)
```

## Tests

The package includes comprehensive test coverage across multiple test files.

### Running Tests
```bash
# Run all tests
go test ./tree

# Run tests with verbose output
go test -v ./tree

# Run specific test
go test -run TestAVLInsert ./tree

# Run tests with race detection
go test -race ./tree

# Run tests with coverage
go test -cover ./tree
```

### Test Coverage
The test suite covers:
- Basic tree operations (insert, delete, search)
- Tree balancing and rotations (AVL-specific)
- Edge cases (empty trees, single nodes, duplicate values)
- Tree traversal in all supported orders
- Tree comparison and equality functions
- ASCII visualization output
- Error handling and boundary conditions

## Benchmarks

The package includes performance benchmarks to compare different tree implementations.

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./tree

# Run specific benchmark
go test -bench=BenchmarkTreeInsert ./tree

# Run benchmarks with memory stats
go test -bench=. -benchmem ./tree

# Run benchmarks multiple times for accuracy
go test -bench=BenchmarkTreeInsert -count=5 ./tree

# Filter benchmarks by tree type
go test -bench=. -args -tree-type=AVL ./tree

# Limit benchmark dataset size
go test -bench=. -args -tree-size-limit=10000 ./tree
```

### Benchmark Flags
The benchmarks support several command-line flags:

- **`-tree-type`**: Filter benchmarks by tree type (BST, AVL)
- **`-tree-size-limit`**: Set maximum number of nodes to insert in a tree for testing
- **`-count`**: Number of times to run the benchmarks

### Example Benchmark Output
```
BenchmarkTreeInsert/BST-0000100-8         50000    25432 ns/op    1024 B/op     12 allocs/op
BenchmarkTreeInsert/AVL-0000100-8         45000    28901 ns/op    1024 B/op     12 allocs/op
BenchmarkTreeInsert/BST-0001000-8          5000   245821 ns/op   10240 B/op    120 allocs/op
BenchmarkTreeInsert/AVL-0001000-8          4800   251234 ns/op   10240 B/op    120 allocs/op
```

## Future Plans

### N-ary Tree Types (Planned)

- **B-Tree**: Self-balancing tree for systems with large branching factors
- **B+-Tree**: B-tree variant optimized for range queries and sequential access

### Additional Expected Features (Planned)
- Tree persistence and serialization
- More comprehensive tree utility functions
- Performance optimizations
- Thread-safe checks and variants
- Custom comparison function support

### Longer term / lower priority features

- Interactive Tree builder taking advantage of the existing methods and rendering capabilities (e.g. Insert, Delete, Search displaying the updated tree after each action)
- Dynamic tree width allocation based on actual tree structure bushy-ness. Right now a 5 level deep tree with 5 digit wide node values ends up being > 160 characters wide regardless of how many nodes are present and where they fall in the overall tree.

## Development

### Contributing

When adding new tree types or features:

1. Implement the `Tree[T]` interface
2. Implement the more specific type of Tree interface (BinaryTree[T] or NAryTree[T])
3. Add comprehensive unit tests covering normal expected values, nil and edge cases, as well as any exceptions that could occur.
4. Include benchmarks comparing with existing implementations
5. Update documentation and examples
6. Ensure all tests pass with `go test ./...`
7. Verify no lint issues with `golangci-lint run ./...`

