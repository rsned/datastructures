# Go Tree Data Structures

This package contains a variety of common ordered, rooted tree types implemented in Go.

There are implementations of a number of traditional Binary Trees (Binary Search Tree,
 AVL Tree, Red/Black tree, etc.) as well as some N-ary trees like B-tree and B+-tree.

The package contains comprehensive unit tests as well as a suite of Benchmarks
that can be used to compare and contrast the various trees performances.

Currently the trees are implemented using generics that support all common Go 
primitive types (integer and floating-point values as well as strings). Future plans
include expanding the library to work with user defined comparable types as well.

R-Tree (and variants) are not currently in scope for this package, but could one
day be added.

## Binary Tree Types


### BST - Binary Search Tree

Simplest of all binary trees, a node, two children, and no balancing.

### AVL - Adelson-Velsky and Landis Tree

A self-balancing binary tree where the height of any two child subtrees differs
by no more than 1.

### Red-Black Tree

An approximately balanced binary tree.

## N-ary Tree Types

In graph theory, an n-ary tree (for nonnegative integers n) is an ordered tree in which each node has no more than n children. 

### B-Tree

A B-tree is a self-balancing tree data structure that maintains sorted data and allows searches, sequential access, insertions, and deletions in logarithmic time. 

### B+-Tree

A B+ tree can be viewed as a B-tree in which each node contains only keys (not key–value pairs),
and to which an additional level is added at the bottom with linked leaves.


## Benchmarks


