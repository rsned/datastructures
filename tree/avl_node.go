package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// avlNode is the actual node in an AVL tree.
type avlNode[T constraints.Ordered] struct {
	value T

	// bf is the balance factor, the height difference of the nodes
	// two subtrees.
	bf int

	// parent is a pointer back to the parent node to allow for updates
	// when rebalancing and navigating.
	parent *avlNode[T]

	// The two children nodes.
	left  *avlNode[T]
	right *avlNode[T]
}

// HasLeft reports if this node has a Left child.
func (t *avlNode[T]) HasLeft() bool {
	return t.left != nil
}

// Left returns this nodes Left child.
func (t *avlNode[T]) Left() BinaryTree[T] {
	return t.left
}

// HasRight reports if this node has a Right child.
func (t *avlNode[T]) HasRight() bool {
	return t.right != nil
}

// Right returns this nodes Right child.
func (t *avlNode[T]) Right() BinaryTree[T] {
	return t.right
}

// Value returns this nodes Value.
func (t *avlNode[T]) Value() T {
	return t.value
}

const (
	balanceFactorSuperscript2    = "²"
	balanceFactorSuperscript1    = "¹"
	balanceFactorSuperscript0    = "⁰"
	balanceFactorSuperscriptNeg1 = "⁻¹"
	balanceFactorSuperscriptNeg2 = "⁻²"
	balanceFactorSubscript2      = "₂"
	balanceFactorSubscript1      = "₁"
	balanceFactorSubscript0      = "₀"
	balanceFactorSubscriptNeg1   = "₋₁"
	balanceFactorSubscriptNeg2   = "₋₂"
)

// Metadata returns a string of metadata about this node.
// For AVL tree, this is the balance factor of the node.
func (t *avlNode[T]) Metadata() string {
	switch t.bf {
	case 2:
		return balanceFactorSubscript2
	case 1:
		return balanceFactorSubscript1
	case 0:
		return balanceFactorSubscript0
	case -1:
		return balanceFactorSubscriptNeg1
	case -2:
		return balanceFactorSubscriptNeg2
	default:
		return fmt.Sprintf("(%d)", t.bf)
	}
}

// balanceFactor returns the nodes balance factor.
// TODO(rsned): Make this public?
func (t *avlNode[T]) balanceFactor() int {
	if t == nil {
		return 0
	}

	return t.right.Height() - t.left.Height()
}

// insertInternal performs AVL insertion with proper rebalancing
// and returns the potential new root.
func (t *avlNode[T]) insertInternal(v T) (*avlNode[T], bool) {
	if t == nil {
		// Create new node
		newNode := &avlNode[T]{
			value:  v,
			bf:     0,
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
	// If we need to go farther left, recurse!
	if v < t.value {
		t.left, inserted = t.left.insertInternal(v)
		if t.left != nil {
			t.left.parent = t
		}
	} else {
		// If we need to go farther right, recurse!
		t.right, inserted = t.right.insertInternal(v)
		if t.right != nil {
			t.right.parent = t
		}
	}

	if !inserted {
		return t, false
	}

	// Update balance factor
	t.bf = t.balanceFactor()

	// Now we need to check for imbalance and apply updates as needed.

	if t.bf > 1 { // The node is right-heavy
		// Check if it's Right-Right or Right-Left
		if t.right != nil && t.right.bf < 0 {
			// Right-Left case
			// Double rotation: Right(Z) then Left(X)
			return rotateRightLeftAVL(t), true
		}
		// Right-Right case
		return rotateLeftAVL(t), true
	} else if t.bf < -1 { // Left-heavy
		// Check if it's Left-Right or Left-Left
		if t.left != nil && t.left.bf > 0 {
			// Left-Right case
			// Double rotation: Left(Z) then Right(X)
			return rotateLeftRightAVL(t), true
		}
		// Left-Left case
		// Single rotation: Right(X)
		return rotateRightAVL(t), true
	}

	return t, true
}

// Insert inserts the node into the tree, growing as needed, and reports
// if the operation was successful.
// NOTE: This method doesn't handle root changes.
// For proper AVL insertion that can change the root, use AVL.Insert() instead.
func (t *avlNode[T]) Insert(v T) bool {
	if t == nil {
		return false
	}

	_, inserted := t.insertInternal(v)

	return inserted
}

// rotateLeftAVL performs a left rotation around the given node.
//
// There are three common forms of transformation:
//
//	parent
//	   \
//	   [H] (+2)
//	     \
//	     [N] (+1)
//	       \
//	       [Z] (0)
//
// Which becomes:
//
//	   parent
//	      \
//	      [N] (0)
//	      / \
//	(0) [H] [Z] (0)
//
// Alternatively, this could be part of a double rotation where we
// are only shifting the node and its child into a form that rotateRight
// will then handle.
//
//	       parent
//	         /
//	  (-2) [C]
//	       /
//	(+1) [A]   <-- node
//	      \
//	  (0) [B]
//
// Which becomes:
//
//	        parent
//	          /
//	   (-2) [C]   <-- node
//	        /
//	 (-1) [B]
//	      /
//	(0) [A]
//
// And finally rotate left with children
//
//	       parent
//	         /
//	       [H] (+2)
//	      /   \
//	(0) [E]   [M] (+1)
//	          / \
//	    (0) [J] [S] (0)
//
// Which becomes:
//
//	       parent
//	          \
//	          [M] (0)
//	         /   \
//	  (0) [H]     [S] (0)
//	      / \
//	(0) [E] [J] (0)
//
// And once again balance is restored.
func rotateLeftAVL[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.right == nil {
		return node
	}

	// Save references
	newRoot := node.right
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

	// Update balance factors after rotation
	node.bf = node.balanceFactor()
	newRoot.bf = newRoot.balanceFactor()

	// Return new root of rotated subtree
	return newRoot
}

// rotateRightAVL performs a right rotation around the given node.
//
// The most common form is:
//
//	       parent
//	          /
//	   (-2) [E]   <-- node
//	        /
//	 (-1) [C]
//	      /
//	(0) [A]
//
// Which becomes:
//
//	   parent
//	      \
//	      [C] (0)   <-- node
//	      / \
//	(0) [A] [E] (0)
//
// Alternatively, this could be part of a double rotation in which case there is
// no grandchild node to handle, we are only shifting shuffling the node and
// its child.
//
//	parent
//	   \
//	   [H] (+2)
//	     \
//	     [Z] (-1)   <-- node
//	     /
//	   [N] (0)
//
// Which becomes:
//
//	parent
//	   \
//	   [H] (_2)   <-- node
//	     \
//	     [N] (+1)
//	       \
//	       [Z] (0)
//
// And now the tree is ready for the rotateLeft to finish the balancing.
//
// The third form is rotate right with children
//
//	            parent
//	              /
//	      (-2)  [H]  <-- node
//	           /   \
//	    (-1) [E]   [J (0)
//	         / \
//	  (-1) [C] [F] (0)
//	      /
//	(0) [A]
//
// Which becomes:
//
//	          parent
//	             /
//	        (0) [E]  <-- node
//	           /   \
//	   (-1) [C]     [H] (0)
//	       /        / \
//	(0) [A] (0)   [F] [J] (0)
//
// And once again balance is restored.
func rotateRightAVL[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.left == nil {
		return node
	}

	// Save references
	newRoot := node.left
	subtree := newRoot.right

	// Perform rotation
	newRoot.right = node
	node.left = subtree

	// Update parent pointers
	newRoot.parent = node.parent
	node.parent = newRoot
	if subtree != nil {
		subtree.parent = node
	}

	// Update balance factors after rotation
	node.bf = node.balanceFactor()
	newRoot.bf = newRoot.balanceFactor()

	return newRoot
}

// rotateRightLeftAVL performs a double rotation: right rotation followed by
// left rotation. This handles the Right-Left case in AVL rebalancing.
//
//	       Step 1:         Step 2:         Result:
//
//		      3               3               5
//			 / \             / \             / \
//			1   7    =>     1   5     =>    3   7
//			   / \             / \         / \ / \
//			  5   8           4   7       1  4 6  8
//			 / \                 / \
//			4   6               6   8
func rotateRightLeftAVL[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.right == nil {
		return node
	}

	// First rotation: right rotation on node.right
	node.right = rotateRightAVL(node.right)
	// Update parent pointer
	if node.right != nil {
		node.right.parent = node
	}

	// Second rotation: left rotation on node
	return rotateLeftAVL(node)
}

// rotateLeftRightAVL performs a double rotation: left rotation
// followed by right rotation.
// This handles the Left-Right case in AVL rebalancing.
//
// Step 1:         Step 2:          Result:
//
//	        7             7               5
//		   / \           / \             / \
//		  3   8   =>    5   8     =>    3   7
//		 / \           / \             / \ / \
//		1   5         3   6           1  4 6  8
//		   / \       / \
//		  4   6     1   4
func rotateLeftRightAVL[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	if node == nil || node.left == nil {
		return node
	}

	// First rotation: left rotation on node.left
	node.left = rotateLeftAVL(node.left)
	// Update parent pointer
	if node.left != nil {
		node.left.parent = node
	}

	// Second rotation: right rotation on node
	return rotateRightAVL(node)
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *avlNode[T]) Delete(_ T) bool {
	if t == nil {
		return false
	}

	return false
}

// Search reports if the given value is in the tree.
func (t *avlNode[T]) Search(v T) bool {
	// If this (child) node is nil, then there is nothing to find.
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

// Traverse traverses the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
//
// NOTE: Nodes in general are not expected to initiate the traverse. It would
// normally be kicked off by the main container type, e.g., AVL not avlNode.
func (t *avlNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)

	// If the node is nil that we are trying to traverse, return the channel,
	// but close it off since there is no way to have anything to send.
	if t == nil {
		defer close(ch)

		return ch
	}

	go func() {
		traverseBinaryTree(t, tOrder, ch)
		close(ch)
	}()

	return ch
}

// Height returns the height of the longest path in the tree from the
// root node to the farthest leaf.
func (t *avlNode[T]) Height() int {
	if t == nil {
		return 0
	}

	lHeight := t.left.Height()
	rHeight := t.right.Height()

	if lHeight > rHeight {
		return lHeight + 1
	}

	return rHeight + 1
}

// toTestString prints out this node with all its properties and children as a
// formatted Go value ready to copy and paste into test code. The parent
// pointer is not set here because it is created when the variable
// is instantiated. The indent param tells how deep in the tree we are so the
// value comes out already gofmt'ed.
func (t *avlNode[T]) toTestString(buf *bytes.Buffer, indent int) {
	// testIndents is a sequence of tab characaters that are to be substringed
	// at the necessary level for proper indenting of node text.
	const testIndents = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"

	fmt.Fprintf(buf, "%svalue: %v,\n", testIndents[:indent], t.value)
	fmt.Fprintf(buf, "%sbf: %d,\n", testIndents[:indent], t.bf)

	if t.left != nil {
		buf.WriteString(testIndents[:indent] + "left: &avlNode[T]{\n")
		t.left.toTestString(buf, indent+1)
		buf.WriteString(testIndents[:indent] + "},\n")
	}

	if t.right != nil {
		buf.WriteString(testIndents[:indent] + "right: &avlNode[T]{\n")
		t.right.toTestString(buf, indent+1)
		buf.WriteString(testIndents[:indent] + "},\n")
	}
}

// Clone creates a deep copy of this AVL node and its subtree.
func (t *avlNode[T]) Clone() Tree[T] {
	if t == nil {
		return nil
	}

	clone := &avlNode[T]{
		value:  t.value,
		bf:     t.bf,
		parent: nil, // Parent will be set during tree construction
		left:   nil,
		right:  nil,
	}

	if t.left != nil {
		if leftClone, ok := t.left.Clone().(*avlNode[T]); ok {
			clone.left = leftClone
			leftClone.parent = clone
		}
	}

	if t.right != nil {
		if rightClone, ok := t.right.Clone().(*avlNode[T]); ok {
			clone.right = rightClone
			rightClone.parent = clone
		}
	}

	return clone
}
