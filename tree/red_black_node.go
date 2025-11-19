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

// insertInternal performs BST insertion and returns the newly inserted leaf node.
// The tree root is NOT changed by this function - fixup must be called separately.
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
		return nil, false
	}

	var inserted bool
	var newNode *redBlackNode[T]
	// If we need to go farther left, recurse.
	if v < t.value {
		newNode, inserted = t.left.insertInternal(v)
		if inserted && t.left == nil {
			// Only update if we're at the insertion point
			t.left = newNode
			newNode.parent = t
		}
	} else {
		// If we need to go farther right, recurse.
		newNode, inserted = t.right.insertInternal(v)
		if inserted && t.right == nil {
			// Only update if we're at the insertion point
			t.right = newNode
			newNode.parent = t
		}
	}

	// Return the newly inserted leaf node
	return newNode, inserted
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

// rotateRightRedBlack performs a right rotation around the given node to
// rebalance but does NOT recolor the nodes.
//
// When the newly inserted node's "uncle" (the sibling of the parent node) is
// red, you can often fix by recoloring. If the uncle is black or nil, you
// need to perform a rotation.
//
// Example: Insert 10 into this tree.
//
//	   [30]         [30]
//	   /            /
//	(20)    ==>   (20)
//	              /
//	            (10)
//
// Parent (20) is red, Uncle is nil (considered black)
// This is a left-left case: 20 is the left child of 30, and 10 is the
// left child of 20
// Fix needed: Right rotation on node 20
//
// The right rotation makes 20 the new root of this subtree, with 10 as its
// left child and 30 becoming its right child:
//
//	   (20)
//	  /    \
//	[10]   (30)
func rotateRightRedBlack[T constraints.Ordered](node *redBlackNode[T]) *redBlackNode[T] {
	if node == nil || node.left == nil {
		return node
	}

	// The new root will be the left pointer.
	newRoot := node.left
	// Save any existing right subtree
	subtree := newRoot.right

	// Perform the rotation
	newRoot.right = node
	node.left = subtree

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

	return newRoot
}

// insertFixup fixes Red-Black tree violations after insertion and returns
// the potential new root of the tree.
//
// There are three potential fixes that can be applied:
//
// - Case 1:  Recoloring
// - Case 2:  Zig-Zag (Triangle) - Left-Right or Right-Left
// - Case 3:  Line-Line
func insertFixup[T constraints.Ordered](node *redBlackNode[T]) *redBlackNode[T] {
	current := node
	for current.parent != nil && current.parent.isRed {
		// Determine the potential uncle node:
		//
		//         G                    G
		//       /   \                /   \
		//      P     U    - or -    U     P
		//     / \   / \            / \   / \
		//    C   S *   *          *   * S   C
		//
		// C is the current node, P is the parent node, U is the uncle node,
		// G is the grandparent node, S is the sibling node, * is a nil/black node.
		//
		// The uncle is the sibling of the parent.
		//
		// The grandparent is the parent of the parent.

		grandparent := current.parent.parent
		if grandparent == nil {
			break
		}

		var uncle *redBlackNode[T]
		if current.parent == grandparent.left {
			uncle = grandparent.right
		} else {
			uncle = grandparent.left
		}

		// Case 1: Recoloring. Uncle is red - recolor to black.
		if uncle != nil && uncle.isRed {
			current.parent.isRed = false
			uncle.isRed = false
			grandparent.isRed = true
			current = grandparent

			continue
		}

		// Cases 2 & 3: Uncle is black (i.e., a leaf node)

		// Determine if we have a zig-zag (triangle) configuration
		//
		//  Left-Right   Right-Left
		//     [G]         [G]
		//     /             \
		//   (P)             (P)
		//     \             /
		//     (C)         (C)
		//
		// If this is Left-??? (parent is left child of grandparent)
		if current.parent == grandparent.left {
			parent := current.parent
			if current == parent.right {
				// Left-Right - rotate parent left to transform to Left-Left
				rotateLeftRedBlack(parent)
				// After rotation, current is now the parent, and parent is its left child
				current, parent = parent, current
			}
			// After Left-Right transform or if already Left-Left:
			// Recolor and rotate grandparent right.
			parent.isRed = false
			grandparent.isRed = true
			rotateRightRedBlack(grandparent)

			break
		}

		// If this is Right-something (parent is right child of grandparent)
		parent := current.parent
		if current == parent.left {
			// Right-Left - rotate parent right to transform to Right-Right
			rotateRightRedBlack(parent)
			// After rotation, current is now the parent, and parent is its right child
			current, parent = parent, current
		}
		// After Right-Left transform or if already Right-Right:
		// Recolor and rotate grandparent left.
		parent.isRed = false
		grandparent.isRed = true
		rotateLeftRedBlack(grandparent)

		break
	}

	// Find and return the root of the entire tree
	root := current
	for root.parent != nil {
		root = root.parent
	}

	// Ensure root is always black
	root.isRed = false

	return root
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
