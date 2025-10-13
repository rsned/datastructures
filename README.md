# datastructures
A collection of data structure libraries and utilities in Go. There already exist many of these on the internet. This collection is mostly for my own edification and refreshing all the old knowledge.

## Tree Data Structures

The main implementation in this repository is the [`tree`](tree/) package, which provides a comprehensive collection of tree data structures implemented using Go generics for any type satisfying `constraints.Ordered` (int, float64, string, etc.)

### Implemented Tree Types

- **Binary Tree** - General binary tree interface providing the `Tree[T]`, `Value`, `Left`, and `Right` concepts.
- **Binary Search Tree (BST)** - Basic unbalanced binary search tree
- **AVL Tree** - Self-balancing height-balanced binary tree  
- **Red-Black Tree** - Approximately balanced tree with red/black node coloring

### Key Features

- All trees implement the `Tree[T]` interface for consistent operations (Insert, Delete, Search, Traverse)
- All of the common traversal orders are supported: in-order, pre-order, post-order, reverse-order, and level-order traversal
- ASCII tree display functionality for debugging and visualization (supports up to height 9 before eliding. (Beyond height 5 though, your display needs to get pretty wide pretty quickly)).
- Extensive test coverage along with a variety of benchmarks for performance comparison
- A number of Utility Functions are already implemented with more planned. Tree comparison, equality and equivalence checking, and other helper functions.

### Quick Example

```go
import "github.com/rsned/datastructures/tree"

// Create an AVL tree
avl := tree.NewAVL[int]()
avl.Insert(10)
avl.Insert(5)
avl.Insert(15)

// Search and traverse
found := avl.Search(5)  // true
for value := range avl.Traverse(tree.TraverseInOrder) {
    fmt.Println(value)  // 5, 10, 15
}
```

For detailed documentation, examples, and API reference, see the [`tree` package documentation](tree/README.md).

## Future Plans

Additional tree types are planned including B-Tree and B+-Tree implementations for n-ary tree support.
