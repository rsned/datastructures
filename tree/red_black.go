package tree

import "golang.org/x/exp/constraints"

const treeTypeRedBlack = "Red-Black"

// RedBlack Tree is a binary search tree where each node has an additional
// attribute: a color, which can be either red or black. By maintaining
// specific coloring rules during insertions and deletions, the tree ensures it
// stays approximately balanced.
//
// The Five Rules
// Every Red-Black tree must satisfy these properties:
//
//   - Every node is colored either red or black
//   - The root node is always black
//   - All leaf nodes (nil children nodes) are black
//   - Red nodes cannot have red children (no two red nodes can be adjacent)
//   - Every path from a node to its descendant nil nodes contains the same
//     number of black nodes (this is called the black-height property)
//
// The longest possible path from root to leaf can be at most twice as long as
// the shortest path, which happens when one path alternates red-black nodes
// while another has only black nodes.
type RedBlack[T constraints.Ordered] struct {
	root *redBlackNode[T]
}

// NewRedBlack returns an empty Red-Black tree ready to use.
func NewRedBlack[T constraints.Ordered]() Tree[T] {
	return &RedBlack[T]{
		root: nil,
	}
}

// Root returns the root node of the tree.
func (t *RedBlack[T]) Root() BinaryTree[T] {
	return t.root
}

// Insert inserts the node into the tree, growing as needed.
func (t *RedBlack[T]) Insert(v T) bool {
	if t.root == nil {
		t.root = &redBlackNode[T]{
			value:  v,
			isRed:  false,
			parent: nil,
			left:   nil,
			right:  nil,
		}

		return true
	}

	newRoot, inserted := t.root.insertInternal(v)
	if inserted {
		t.root = newRoot

		/*
			// TODO(rsned): Uncomment this when its working.
			// Apply fixup if there's a red-red violation.
			// Find the newly inserted node to check for violations
			newNode := findNodeRedBlack(t.root, v)

			if newNode != nil && newNode.isRed && newNode.parent != nil && newNode.parent.isRed {
				// We have a red-red violation, fix it up
				// t.root = insertFixup(newNode)
			}
		*/
		// Ensure root is black
		t.root.isRed = false
	}

	return inserted
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
//
// The trees internal structure may be updated.
func (t *RedBlack[T]) Delete(v T) bool {
	if t.root == nil {
		return false
	}

	return t.root.Delete(v)
}

// Search reports if the given value is in the tree.
func (t *RedBlack[T]) Search(v T) bool {
	if t.root == nil {
		return false
	}

	return t.root.Search(v)
}

// Traverse traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *RedBlack[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t.root, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *RedBlack[T]) Height() int {
	if t.root == nil {
		return 0
	}

	return t.root.Height()
}

// Clone creates a deep copy of the Red-Black tree.
func (t *RedBlack[T]) Clone() Tree[T] {
	if t == nil || t.root == nil {
		return NewRedBlack[T]()
	}

	clone := &RedBlack[T]{
		root: nil,
	}
	if clonedRoot, ok := t.root.Clone().(*redBlackNode[T]); ok {
		clone.root = clonedRoot
	}

	return clone
}
