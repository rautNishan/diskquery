package tree

import (
	"fmt"
	"math"
)

type BsNode struct {
	value  int
	left   *BsNode
	right  *BsNode
	height int
}

type BinarySearchTree struct {
	root *BsNode
}

func NewNode(value int) *BsNode {
	return &BsNode{
		value:  value,
		left:   nil,
		right:  nil,
		height: 0,
	}
}

func (bst *BinarySearchTree) IsEmpty() bool {
	return bst.root == nil
}

func (bst *BinarySearchTree) GetHeight(node *BsNode) int {
	if node == nil {
		return -1
	}
	return node.height
}

func (bst *BinarySearchTree) Insert(value int) {
	bst.root = bst.insert(bst.root, value)
}

func (bst *BinarySearchTree) insert(node *BsNode, value int) *BsNode {
	if node == nil {
		newNode := NewNode(value)
		return newNode
	}

	if value < node.value {
		node.left = bst.insert(node.left, value)
	} else {
		node.right = bst.insert(node.right, value)
	}

	node.height = max(bst.GetHeight(node.left), bst.GetHeight(node.right)) + 1

	return node
}

func (bst *BinarySearchTree) IsBalancedTree() bool {
	return bst.isBalanced(bst.root)
}

func (bst *BinarySearchTree) isBalanced(node *BsNode) bool {
	if node == nil {
		return true
	}
	return math.Abs(float64(bst.GetHeight(node.left)-bst.GetHeight(node.right))) <= 1 && bst.isBalanced(node.left) && bst.isBalanced(node.right)
}

func (bst *BinarySearchTree) PrettyPrint() {
	printHelper(bst.root, "", true)
}

func printHelper(node *BsNode, prefix string, isTail bool) {
	if node == nil {
		return
	}

	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		printHelper(node.right, newPrefix, false)
	}

	// Print current node
	fmt.Printf("%s", prefix)
	if isTail {
		fmt.Printf("└── ")
	} else {
		fmt.Printf("┌── ")
	}
	fmt.Printf("%d(h=%d)\n", node.value, node.height)

	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		printHelper(node.left, newPrefix, true)
	}
}
