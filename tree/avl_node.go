package tree

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

// avlNode is the actual node in an AVL tree.
type avlNode[T constraints.Ordered] struct {
	value T

	// bf is the balance factor, the height difference of the nodes two subtrees.
	// Could probably be an int8 since its always in the range [-2, +2]
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

// Metadata returns a string of metadata about this node.
// For AVL tree, this is tghe balance factor of the node.
func (t *avlNode[T]) Metadata() string {
	return fmt.Sprintf("BF:%2d", t.bf)
}

// balanceFactor returns the nodes balance factor.
// TODO(rsned): Make this public?
func (t *avlNode[T]) balanceFactor() int {
	if t == nil {
		return 0
	}
	return t.right.Height() - t.left.Height()
}

// Insert inserts the node into the tree, growing as needed, and reports
// if the operation was successful.
func (t *avlNode[T]) Insert(v T) bool {
	if t == nil {
		t = &avlNode[T]{
			value: v,
			bf:    0,
		}
		return true
	}

	// Inserting a duplicate value is an error.
	if v == t.value {
		return false
	}

	// If we need to go farther left, recurse!
	if v < t.value && t.left != nil {
		return t.left.Insert(v)
	}

	// If we need to go farther right, recurse!
	if v > t.value && t.right != nil {
		return t.right.Insert(v)
	}

	// We are at the end of the line going left, we need to add
	// a new node to the left, updating balancing factors.
	if v < t.value {
		t.left = &avlNode[T]{
			parent: t,
			value:  v,
		}
	} else {
		// Or we need to add a new node to the right.
		t.right = &avlNode[T]{
			parent: t,
			value:  v,
		}
	}

	// Update the balance factor back up from here after adding the new node.
	updateBalanceFactors(t)

	// Now we need to check for imbalance and apply updates as needed.

	// We are in a leaf node, or a node with only the new child, so this node
	// is not the one that needs a rebalance. Start working our way up the tree
	// looking for a node that is far enough out of balance.
	for x := t; x != nil; x = x.parent {
		// If this next level is balanced, move up and try again.
		// We will either get to the root or find an imbalance.
		if x.bf == 0 {
			continue
		}

		if x.bf > 1 { // The node is right-heavy
			if x.right != nil {
				// Check if it's Right-Right or Right-Left
				if x.right.bf < 0 {
					// Right Left Case
					// Double rotation: Right(Z) then Left(X)
					rotateRightLeft(x)
				} else if x.right.bf > 0 {
					// Right Right Case
					// Single rotation Left(X)
					rotateLeft(x)
				}
			}
		} else if x.bf < -1 {
			if x.left != nil {
				// Check if it's Left-Right or Left-Left
				if x.left.bf > 0 {
					// Left Right Case
					// Double rotation: Left(Z) then Right(X)
					rotateLeftRight(x)
				} else if x.left.bf < 0 {
					// Left Left Case
					// Single rotation Right
					rotateRight(x)
				}
			}
		}
	}
	return true
}

func updateBalanceFactors[T constraints.Ordered](node *avlNode[T]) {
	const limit = 4
	var i int
	// Update the balance factor back up from here after adding the new node.
	for x := node; x != nil; x = x.parent {
		x.bf = x.balanceFactor()
		i++
		if i > limit {
			break
		}
	}
}

// rotateLeft takes a node in the tree and rotates left through the middle
// node to balance it.
//
// THe most common form is:
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
// And now the tree has regained balance.
//
// Alternatively, this could be part of a double rotation in which case there is
// no grandchild node to handle, we are only shifting the node and its child into
// a form that rotateRight will then handle.
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
// And now the tree is ready for the rotateRight to finish the balancing.
//
// The third form is a rotate left with children:
//
//	       parent
//	         /
//	       [H] (+1)   <-- node
//	      /   \
//	(0) [E]   [M] (+2)
//	          / \
//	    (0) [J] [S] (+1)
//	              \
//	              [Z] (0)
//
// Which becomes:
//
//	       parent
//	          \
//	          [M] (0)   <-- node
//	         /   \
//	  (0) [H]     [S] (+1)
//	      / \        \
//	(0) [E] [J] (0)   [Z] (0)
//
// And once again balance is restored.
func rotateLeft[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	// node is the node with a balance factor >= 2
	// Save its two children and its right childs two children.
	childL := node.left
	childR := node.right
	grandchildL := childR.left
	grandchildR := childR.right

	// parent
	//   \
	//   [H]  <-- node
	//     \
	//     [N]
	//       \
	//       [Z]
	//

	// Move N to H's left child.
	// Move Z up to H's right spot and update its parent to H.
	//
	// parent
	//   \
	//   [H]
	//   / \
	// [N] [Z]
	//
	node.left = childR
	node.right = grandchildR

	// If this was a full rotate (and not the first part of a rotate left then right)
	// then there would be a grandchild node that would need its parent set.
	if node.right != nil {
		node.right.parent = node
	}

	// If the right child had a left grandchild tree, it jumps over to become
	// the new left childs left node.
	node.left.left = grandchildL
	if node.left.left != nil {
		node.left.left.parent = node.left
	}

	// If there was an existing left child it becomes the left nodes right grandchild.
	node.left.right = childL
	if node.left.right != nil {
		node.left.right.parent = node.left
	}

	// Swap H & N's values
	//
	//  parent
	//    \
	//    [N]
	//    / \
	//  [H] [Z]
	//
	node.value, childR.value = childR.value, node.value

	// Update the affected nodes balance factors and up the tree.
	updateBalanceFactors(node.left)
	// For the other child node, only need to update it by itself.
	// updateBalanceFactors handles the main node and on up.
	if node.right != nil {
		node.right.bf = node.right.balanceFactor()
	}

	// Return new root of rotated subtree
	return node
}

// rotateRight takes a set of nodes and rotates right through the middle
// node to balance it.
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
// no grandchild node to handle, we are only shifting shuffling the node and its child.
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
func rotateRight[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	// Save its two children and its right childs two children.
	childL := node.left
	childR := node.right
	grandchildL := childL.left
	grandchildR := childL.right

	// From our starting point:
	//
	//       parent
	//        /
	//      [E]  (<--node)
	//      /
	//    [C]
	//    /
	//  [A]
	//
	// Move left child to node's right.
	// Move left grandchild up to left child and update its parent to node..
	//
	//   parent
	//     \
	//     [E]
	//     / \
	//   [A] [C]
	//
	node.left = grandchildL
	node.right = childL
	// If this was a full rotate (and not the first part of a rotate right then left)
	// then there would be a grandchild node that would need its parent set.
	if node.left != nil {
		node.left.parent = node
	}

	// If the left child had a right grandchild tree, it jumps over to become
	// the new right childs left node.
	node.right.left = grandchildR
	if node.right.left != nil {
		node.right.left.parent = node.right
	}

	// If there was an existing right child it becomes nodes right grandchild.
	node.right.right = childR
	if node.right.right != nil {
		node.right.right.parent = node.right
	}

	// Swap node and new right childs values.
	//
	//  parent
	//    \
	//    [C]
	//    / \
	//  [A] [E]
	//
	node.value, childL.value = childL.value, node.value

	// Update the affected nodes balance factors and the parents.
	updateBalanceFactors(node.left)
	// For the other child node, only need to update it by itself.
	// updateBalanceFactors handles the main node and on up.
	if node.right != nil {
		node.right.bf = node.right.balanceFactor()
	}

	// Return new root of rotated subtree
	return node
}

// rotateRightLeft performs a double rotation, first right around the middle node
// to transform it into the standard form for the follow up rotateLeft.
//
//	 \
//	[ H ] (+2)
//	    \
//	   [ N ] (-1)
//	   /
//	[ K ] (0)
//
// becomes:
//
//	 \
//	[ H ] (+2)
//	    \
//	   [ K ] (+1)
//	      \
//	      [ N ] (0)
//
// which becomes:
//
//	    \
//	   [ K ] (0)
//	   /   \
//	[ H ] [ Z ]
//	 (0)   (0)
//
// And balance is once again restored.
func rotateRightLeft[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	rotateRight(node.right)
	rotateLeft(node)
	return node
}

func rotateLeftRight[T constraints.Ordered](node *avlNode[T]) *avlNode[T] {
	rotateLeft(node.left)
	rotateRight(node)
	return node
}

// Delete the requested node from the tree and reports if it was successful.
// If the value is not in the tree, the tree is unchanged and false is returned.
// If the node is not a leaf the trees internal structure may be updated.
func (t *avlNode[T]) Delete(v T) bool {
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

// Traverses traverse the tree in the specified order emitting the values to
// the channel. Channel is closed once the final value is emitted.
//
// NOTE: Nodes in general are not expected to initiate the traverse. It would
// normally be kicked off by the main container type, e.g., AVL not avlNode.
func (t *avlNode[T]) Traverse(tOrder TraverseOrder) <-chan T {
	ch := make(chan T)

	// If the node is nil that we are trying to traverse, return the channel,
	// but close it off since there is no way to have anythign to send.
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
	lh := t.left.Height()
	rh := t.right.Height()

	if lh > rh {
		return lh + 1
	}

	return rh + 1
}

func (t *avlNode[T]) toTestString(buf *bytes.Buffer, indent int) {
	buf.WriteString(fmt.Sprintf("%svalue: %v,\n", testIndents[:indent], t.value))
	buf.WriteString(fmt.Sprintf("%sbf: %d,\n", testIndents[:indent], t.bf))

	if t.left != nil {
		buf.WriteString(fmt.Sprintf("%sleft: &avlNode[T]{\n", testIndents[:indent]))
		t.left.toTestString(buf, indent+1)
		buf.WriteString(fmt.Sprintf("%s},\n", testIndents[:indent]))
	}
	if t.right != nil {
		buf.WriteString(fmt.Sprintf("%sright: &avlNode[T]{\n", testIndents[:indent]))
		t.right.toTestString(buf, indent+1)
		buf.WriteString(fmt.Sprintf("%s},\n", testIndents[:indent]))
	}
}
