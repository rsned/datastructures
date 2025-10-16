# Go Data Structures

[![Go Reference](https://pkg.go.dev/badge/github.com/rsned/datastructures.svg)](https://pkg.go.dev/github.com/rsned/datastructures)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsned/datastructures)](https://goreportcard.com/report/github.com/rsned/datastructures)

A collection of well-documented, generic data structure implementations in Go. This repository is intended for both educational purposes and practical use, providing clear and reusable components.

## Overview

This project provides robust implementations of common data structures, built with Go generics to be type-safe and flexible. The primary focus is currently on tree-based structures, with plans to expand to other types in the future.

The main package in this repository is [`tree`](./tree/), which offers a comprehensive collection of tree data structures.

## Key Features

*   **Generic and Type-Safe**: Utilizes Go generics to work with any type that satisfies `constraints.Ordered` (e.g., `int`, `float64`, `string`).
*   **Consistent API**: All tree types implement a common `Tree[T]` interface for standard operations like `Insert`, `Delete`, `Search`, and `Traverse`.
*   **Multiple Tree Implementations**:
    *   **Binary Search Tree (BST)**: A basic, unbalanced binary search tree.
    *   **AVL Tree**: A self-balancing tree that guarantees O(log n) performance for key operations.
    *   **Red-Black Tree**: Another self-balancing tree offering a good trade-off between insertion and search performance.
*   **Rich Functionality**:
    *   Supports all standard traversal orders (In-Order, Pre-Order, Post-Order, etc.).
    *   Includes utility functions for tree comparison (`Equal`, `Equivalent`), conversion, and analysis.
    *   Provides ASCII visualization for debugging and display.
*   **Thoroughly Tested and Documented**: Comes with extensive unit tests, benchmarks, and complete GoDoc documentation.

## Getting Started

To use this library in your project, you can add it with `go get`:

```sh
go get github.com/rsned/datastructures
```

### Quick Example

Here's a simple example of how to create and use an AVL tree:

```go
package main

import (
	"fmt"
	"github.com/rsned/datastructures/tree"
)

func main() {
	// Create a new AVL tree for integers.
	avl := tree.NewAVL[int]()

	// Insert some values.
	avl.Insert(10)
	avl.Insert(5)
	avl.Insert(15)
	avl.Insert(3)
	avl.Insert(7)

	// Search for a value.
	if avl.Search(7) {
		fmt.Println("Found 7 in the tree.")
	}

	// Traverse the tree in-order to get sorted values.
	fmt.Println("In-order traversal:")
	for value := range avl.Traverse(tree.TraverseInOrder) {
		fmt.Printf("%d ", value)
	}
	fmt.Println()

	// Print a visual representation of the tree.
	if root, ok := avl.Root().(tree.BinaryTree[int]); ok {
		fmt.Println("\nASCII Visualization:")
		fmt.Println(tree.PrintBinaryTreeASCII("", root))
	}
}
```

## Documentation

For detailed documentation, examples, and API references, please see the Go documentation or the specific README files within the sub-packages:

*   **[Package Documentation](https://pkg.go.dev/github.com/rsned/datastructures)**
*   **[`tree` Package README](./tree/README.md)**

## Future Plans

*   Implement deletion for AVL and Red-Black trees.
*   Add B-Tree and B+-Tree implementations for n-ary support.
*   Expand the collection to include other data structures like graphs, heaps, and hash maps.