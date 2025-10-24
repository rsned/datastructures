package tree

import "golang.org/x/exp/constraints"

// redBlackNode is the basic node in a Red-Black binary search tree.
type redBlackNode[T constraints.Ordered] struct {
	value T

	isRed bool

	parent      *redBlackNode[T]
	left, right *redBlackNode[T]
}

// HasLeft reports if this node has a Left child.
func (t *redBlackNode[T]) HasLeft() bool {
	return t.left != nil
}

// HasRight reports if this node has a Right child.
func (t *redBlackNode[T]) HasRight() bool {
	return t.right != nil
}

// Left returns this nodes Left child.
func (t *redBlackNode[T]) Left() BinaryTree[T] {
	return t.left
}

// Right returns this nodes Right child.
func (t *redBlackNode[T]) Right() BinaryTree[T] {
	return t.right
}

// Value returns this nodes Value.
func (t *redBlackNode[T]) Value() T {
	return t.value
}

// Metadata returns a string of metadata about this node. In this case if
// the node is a Red or Black node.
func (t *redBlackNode[T]) Metadata() string {
	if t.isRed {
		return "Red"
	}

	return "Black"
}

// Insert inserts the node into the tree, growing as needed, and reports
// if the operation was successful.
// NOTE: This method doesn't handle root changes.
// For proper Red-Black insertion that can change the root, use RedBlack.Insert() instead.
func (t *redBlackNode[T]) Insert(v T) bool {
	if t == nil {
		return false
	}

	_, inserted := t.insertInternal(v)

	return inserted
}

// insertInternal performs Red-Black insertion with proper rebalancing
// and returns the potential new root.
func (t *redBlackNode[T]) insertInternal(v T) (*redBlackNode[T], bool) {
	// We are either at the end of the road drilling down to the right spot in
	// the tree, or we are inserting what will be the new root node.
	if t == nil {
		// New nodes all start as Red
		newNode := &redBlackNode[T]{
			value:  v,
			isRed:  true,
			parent: nil,
			left:   nil,
			right:  nil,
		}

		return newNode, true
	}

	// Inserting a duplicate value is an error
	if v == t.value {
		return t, false
	}

	var inserted bool
	var newNode *redBlackNode[T]
	// If we need to go farther left, recurse.
	if v < t.value {
		newNode, inserted = t.left.insertInternal(v)
		if newNode != nil {
			t.left = newNode
			newNode.parent = t
		}
	} else {
		// If we need to go farther right, recurse.
		newNode, inserted = t.right.insertInternal(v)
		if newNode != nil {
			t.right = newNode
			newNode.parent = t
		}
	}

	if !inserted {
		return t, false
	}

	// TODO(rsned): Uncomment this when its working.
	/*
		if newNode != nil && newNode.isRed && newNode.parent != nil && newNode.parent.isRed {
			// We have a red-red violation, fix it up starting from the new node

			// TODO(rsned): Uncomment this when its working.
			// return insertFixup(newNode), true
		}
	*/

	return t, true
}

// rotateLeftRedBlack performs a left rotation around the given node to
// rebalance but does NOT recolor the nodes.
//
// When the newly inserted node's "uncle" (the sibling of the parent node) is
// red, you can often fix by recoloring. If the uncle is black or nil, you
// need to perform a rotation.
//
// Example: Insert 30 into this tree.
//
//	[10]        [10]
//	  \           \
//	  (20)  ==>   (20)
//	                \
//	                (30)
//
// Parent (20) is red, Uncle is nil (considered black)
// This is a right-right case: 20 is the right child of 10, and 30 is the
// right child of 20
// Fix needed: Left rotation on node 10
//
// The left rotation makes 20 the new root of this subtree, with 10
// becoming its left child:
//
//	   (20)
//	  /    \
//	[10]   (30)
func rotateLeftRedBlack[T constraints.Ordered](node *redBlackNode[T]) *redBlackNode[T] {
	if node == nil || node.right == nil {
		return node
	}

	// The new root will be the right pointer.
	newRoot := node.right
	// Save any existing left subtree
	subtree := newRoot.left

	// Perform rotation
	newRoot.left = node
	node.right = subtree

	// Update parent pointers
	newRoot.parent = node.parent
	node.parent = newRoot
	if subtree != nil {
		subtree.parent = node
	}

	// Update parent's child pointer if it exists
	if newRoot.parent != nil {
		if newRoot.parent.left == node {
			newRoot.parent.left = newRoot
		} else {
			newRoot.parent.right = newRoot
		}
	}

	// Return new root of rotated subtree
	return newRoot
}

// findNodeRedBlack searches for a node with the given value in the tree.
// Returns the node's pointer if found, or nil otherwise.
func findNodeRedBlack[T constraints.Ordered](root *redBlackNode[T], value T) *redBlackNode[T] {
	if root == nil {
		return nil
	}

	if value == root.value {
		return root
	}

	if value < root.value {
		return findNodeRedBlack(root.left, value)
	}

	return findNodeRedBlack(root.right, value)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *redBlackNode[T]) Delete(_ T) bool {
	if t == nil {
		return false
	}

	return false
}

// Search reports if the given value is in the tree.
func (t *redBlackNode[T]) Search(v T) bool {
	if t == nil {
		return false
	}
	if v == t.value {
		return true
	}

	if v < t.value {
		return t.left.Search(v)
	}

	return t.right.Search(v)
}

// Walk traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
func (t *redBlackNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)
	go func() {
		traverseBinaryTree(t, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *redBlackNode[T]) Height() int {
	if t == nil {
		return 0
	}
	lh := t.left.Height()
	rh := t.right.Height()
	if lh > rh {
		return lh + 1
	}

	return rh + 1
}

// Clone creates a deep copy of this Red-Black node and its subtree.
func (t *redBlackNode[T]) Clone() Tree[T] {
	if t == nil {
		return nil
	}

	clone := &redBlackNode[T]{
		value:  t.value,
		isRed:  t.isRed,
		parent: t.parent,
		left:   nil,
		right:  nil,
	}

	if t.left != nil {
		if leftClone, ok := t.left.Clone().(*redBlackNode[T]); ok {
			clone.left = leftClone
		}
	}

	if t.right != nil {
		if rightClone, ok := t.right.Clone().(*redBlackNode[T]); ok {
			clone.right = rightClone
		}
	}

	return clone
}
