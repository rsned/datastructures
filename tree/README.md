# Tree Data Structures

This package provides a collection of generic, type-safe tree data structures implemented in Go. It supports any data type that satisfies `constraints.Ordered`, such as integers, floats, and strings. The implementations are designed to be clear, well-documented, and suitable for both educational and practical use.

## Features

*   **Generic Implementation**: Works with any ordered type thanks to Go generics.
*   **Multiple Tree Types**:
    *   **Binary Search Tree (BST)**: A simple, unbalanced binary tree.
    *   **AVL Tree**: A self-balancing binary search tree that maintains its height for O(log n) performance.
    *   **Red-Black Tree**: A self-balancing binary search tree that uses node coloring to ensure approximate balance.
*   **Consistent API**: All trees adhere to the `Tree[T]` interface for uniform operations (`Insert`, `Delete`, `Search`, `Traverse`).
*   **Standard Traversals**: Supports in-order, pre-order, post-order, reverse-order, and level-order traversals.
*   **Utility Functions**: Includes helpers for comparing trees (`Equal`, `Equivalent`), converting between types, and more.
*   **Visualization**: Provides an ASCII-based tree renderer (`PrintBinaryTreeASCII`) for debugging and display.
*   **Comprehensive Testing**: Includes extensive unit tests and performance benchmarks.

## Core Interfaces

### `Tree[T]` Interface

This is the main interface that all tree types in the package implement.

```go
type Tree[T constraints.Ordered] interface {
    Insert(v T) bool
    Delete(v T) bool
    Search(v T) bool
    Height() int
    Clone() Tree[T]
    Traverse(order TraverseOrder) <-chan T
}
```

### `BinaryTree[T]` Interface

This interface is specific to binary tree nodes and extends the base `Tree[T]` interface.

```go
type BinaryTree[T constraints.Ordered] interface {
    Tree[T]
    Value() T
    Left() BinaryTree[T]
    Right() BinaryTree[T]
    HasLeft() bool
    HasRight() bool
    Metadata() string // Returns node-specific info like balance factor or color.
}
```

## Usage Example

### Creating and Using a Tree

```go
package main

import (
	"fmt"
	"github.com/rsned/datastructures/tree"
)

func main() {
	// Create a new AVL tree for integers.
	avl := tree.NewAVL[int]()

	// Insert values. The tree will self-balance.
	avl.Insert(20)
	avl.Insert(10)
	avl.Insert(30)
	avl.Insert(5)
	avl.Insert(15)

	// Search for a value.
	fmt.Printf("Search for 15: %t\n", avl.Search(15)) // true

	// Traverse the tree to get sorted values.
	fmt.Print("In-order traversal: ")
	for v := range avl.Traverse(tree.TraverseInOrder) {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}
```

### Visualizing a Tree

You can print a tree's structure to the console using `PrintBinaryTreeASCII`.

```go
// Assuming 'avl' is the tree from the example above.
if root, ok := avl.Root().(tree.BinaryTree[int]); ok {
    fmt.Println(tree.PrintBinaryTreeASCII("AVL Tree Structure:", root))
}
```

**Output:**

```
AVL Tree Structure:

     20
    (0)
    /  \
  10    30
 (0)   (0)
 / \
5   15
(0) (0)
```

## Running Tests and Benchmarks

This package comes with a full suite of tests and benchmarks.

### Running Tests

To run all unit tests for the `tree` package:

```sh
go test -v github.com/rsned/datastructures/tree
```

To run tests with code coverage:

```sh
go test -cover github.com/rsned/datastructures/tree
```

### Running Benchmarks

To run all benchmarks and see performance comparisons between tree types:

```sh
go test -bench=. github.com/rsned/datastructures/tree
```

To include memory allocation statistics:

```sh
go test -bench=. -benchmem github.com/rsned/datastructures/tree
```

## Future Development

*   Complete the implementation of `Delete` for AVL and Red-Black trees.
*   Add implementations for n-ary trees like B-Trees and B+-Trees.
*   Introduce more utility functions for advanced tree manipulation.

## Contributing

Contributions are welcome! If you'd like to contribute, please:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Add or update tests for your changes.
4.  Ensure all tests pass (`go test ./...`).
5.  Run `golangci-lint run ./...` to check for linting issues.
6.  Submit a pull request.