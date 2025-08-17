package tree

import "fmt"

type AVLNode struct {
	left   *AVLNode
	right  *AVLNode
	height int
	value  int
}

func (avl *AVLTree) GetHeight() int {
	return avl.getHeight(avl.root)
}

func (avl *AVLTree) getHeight(node *AVLNode) int {
	if node == nil {
		return -1
	}
	return node.height

}

type AVLTree struct {
	root *AVLNode
}

func (avl *AVLTree) Insert(value int) {
	avl.root = avl.insert(avl.root, value)
}

func (avl *AVLTree) insert(node *AVLNode, value int) *AVLNode {
	if node == nil {
		return &AVLNode{
			value:  value,
			left:   nil,
			right:  nil,
			height: 0,
		}
	}

	//Insert Right
	if value > node.value {
		node.right = avl.insert(node.right, value)
	} else {
		//Insert left
		node.left = avl.insert(node.left, value)
	}
	node.height = max(avl.getHeight(node.left), avl.getHeight(node.right)) + 1
	return avl.rotate(node)
}

func (avl *AVLTree) rotate(node *AVLNode) *AVLNode {
	//For Left Heavy
	if avl.getHeight(node.left)-avl.getHeight(node.right) > 1 {
		if avl.getHeight(node.left.left)-avl.getHeight(node.left.right) < 0 {
			node.left = avl.leftRotate(node.left)
			return avl.rightRotate(node)
		} else {
			return avl.rightRotate(node)
		}
	}

	//For Right heavy
	if avl.getHeight(node.left)-avl.getHeight(node.right) < -1 {
		if avl.getHeight(node.right.left)-avl.getHeight(node.right.right) > 0 {
			node.right = avl.rightRotate(node.right)
			return avl.leftRotate(node)
		} else {
			return avl.leftRotate(node)
		}
	}

	return node
}

func (avl *AVLTree) rightRotate(node *AVLNode) *AVLNode {
	b := node.left
	c := b.right

	b.right = node
	node.left = c
	node.height = max(avl.getHeight(node.left), avl.getHeight(node.right)) + 1

	b.height = max(avl.getHeight(b.left), avl.getHeight(b.right)) + 1

	return b
}

func (avl *AVLTree) leftRotate(node *AVLNode) *AVLNode {
	b := node.right
	c := b.left

	b.left = node
	node.right = c
	node.height = max(avl.getHeight(node.left), avl.getHeight(node.right)) + 1

	b.height = max(avl.getHeight(b.left), avl.getHeight(b.right)) + 1

	return b
}

func (avl *AVLTree) PrettyPrint() {
	if avl.root == nil {
		fmt.Println("Empty tree")
		return
	}
	avl.prettyPrint(avl.root, "", false)
}

// prettyPrint is the recursive helper function
func (avl *AVLTree) prettyPrint(node *AVLNode, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	// Print the right subtree first (so it appears at the top)
	if node.right != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		avl.prettyPrint(node.right, newPrefix, false)
	}

	// Print the current node
	connector := "├── "
	if isLeft {
		connector = "└── "
	}
	fmt.Printf("%s%s%d (h:%d)\n", prefix, connector, node.value, node.height)

	// Print the left subtree
	if node.left != nil {
		newPrefix := prefix
		if isLeft {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		avl.prettyPrint(node.left, newPrefix, true)
	}
}

func (avl *AVLTree) IsBalanced() bool {
	return avl.isBalance(avl.root)
}

func (avl *AVLTree) isBalance(node *AVLNode) bool {
	if node == nil {
		return true
	}

	balance := avl.getHeight(node.left) - avl.getHeight(node.right)
	if abs(balance) > 1 {
		return false
	}

	return avl.isBalance(node.left) && avl.isBalance(node.right)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
