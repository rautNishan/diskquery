package tree

import "fmt"

type Node struct {
	value int
	right *Node
	left  *Node
}

func newNode(val int) *Node {
	return &Node{
		value: val,
		left:  nil,
		right: nil,
	}
}

type BinaryTree struct {
	root *Node
}

func (bt *BinaryTree) Populate() {
	fmt.Println("Enter value for root node: ")
	var value int

	_, err := fmt.Scan(&value)
	if err != nil {
		fmt.Println("Error while taking input: ", err)
	}
	bt.root = newNode(value)
	populate(bt.root)
}

func populate(node *Node) {
	fmt.Println("Do you want to insert in left of node ", node.value)
	var leftInsert bool
	_, err := fmt.Scan(&leftInsert)
	if err != nil {
		fmt.Println("Error while taking input: ", err)
	}

	if leftInsert {
		fmt.Println("Enter value for root node: ")
		var value int
		_, err := fmt.Scan(&value)
		if err != nil {
			fmt.Println("Error while taking input: ", err)
		}
		node.left = newNode(value)
		populate(node.left)
	}

	fmt.Println("Do you want to insert in right of node ", node.value)
	var rightInsert bool
	_, err = fmt.Scan(&rightInsert)
	if err != nil {
		fmt.Println("Error while taking input: ", err)
	}
	if rightInsert {
		fmt.Println("Enter value for root node: ")
		var value int
		_, err := fmt.Scan(&value)
		if err != nil {
			fmt.Println("Error while taking input: ", err)
		}
		node.right = newNode(value)
		populate(node.right)
	}
}

func (bt *BinaryTree) PrettyPrint() {
	fmt.Println("**********************PRINTING**********************")
	printTree(bt.root, "", true)
}

func printTree(node *Node, prefix string, isTail bool) {
	if node == nil {
		return
	}

	// Print the current node
	connector := "└── "
	if !isTail {
		connector = "├── "
	}
	fmt.Println(prefix + connector + fmt.Sprintf("%d", node.value))

	// Calculate new prefix for children
	newPrefix := prefix
	if isTail {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	// Count children
	children := []*Node{}
	if node.left != nil {
		children = append(children, node.left)
	}
	if node.right != nil {
		children = append(children, node.right)
	}

	// Recurse
	for i, child := range children {
		printTree(child, newPrefix, i == len(children)-1)
	}
}
