package tree

type AVLNode struct {
	left   *AVLNode
	right  *AVLNode
	height int
	value  int
}

func (node *AVLNode) newNode(value int) *AVLNode {
	return &AVLNode{
		value:  value,
		left:   nil,
		right:  nil,
		height: 0,
	}
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
		newNode := node.newNode(value)
		return newNode
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

	}
	return node
}
