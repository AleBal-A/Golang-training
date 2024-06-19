package main

import "fmt"

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func areToysBalanced(root *TreeNode) bool {
	if root != nil {
		left := countToys(root.Left)
		right := countToys(root.Right)
		return left == right
	}
	return false
}

func countToys(root *TreeNode) int {
	if root == nil {
		return 0
	}
	toys := 0
	if root.HasToy {
		toys = 1
	}
	toys += countToys(root.Left)
	toys += countToys(root.Right)

	return toys
}

func main() {
	tree := &TreeNode{
		HasToy: true,
		Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
		Right:  &TreeNode{HasToy: false, Left: nil, Right: nil},
	}

	fmt.Println(areToysBalanced(tree))
}
