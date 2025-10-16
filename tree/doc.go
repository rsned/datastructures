// Package tree implements a collection of generic, type-safe tree data
// structures. It provides a variety of common binary tree types, with plans to
// expand to n-ary trees in the future.
//
// All tree implementations are built using Go generics, allowing them to work
// with any data type that satisfies the `golang.org/x/exp/constraints.Ordered`
// interface (e.g., integers, floats, strings).
//
// Key Features:
//
//   - A consistent `Tree[T]` interface for common operations like Insert,
//     Delete, Search, and Traverse.
//   - Support for standard traversal orders: in-order, pre-order, post-order,
//     reverse-order, and level-order.
//   - Implementations of several self-balancing trees, including AVL and
//     Red-Black trees, to ensure efficient performance.
//   - Utility functions for tree manipulation, comparison, and visualization.
//
// Available Tree Types:
//
//   - Binary Search Tree (BST): A basic, unbalanced binary search tree.
//   - AVL Tree: A self-balancing binary search tree that maintains its height
//     to provide O(log n) time complexity for search, insert, and delete
//     operations.
//   - Red-Black Tree: Another self-balancing binary search tree that uses
//     node coloring to ensure the tree remains approximately balanced.
package tree