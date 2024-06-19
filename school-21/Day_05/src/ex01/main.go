package main

import (
	"fmt"
)

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func unrollGarland(root *TreeNode) []bool {
	if root == nil {
		return []bool{}
	}
	var result []bool
	queue := []*TreeNode{root}
	leftToRight := false

	for len(queue) > 0 {
		levelSize := len(queue)
		level := make([]bool, levelSize)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			if leftToRight {
				level[i] = node.HasToy
			} else {
				level[levelSize-1-i] = node.HasToy
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level...)
		leftToRight = !leftToRight
	}

	return result
}

func main() {
	tree1 := &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: false,
			Left:   &TreeNode{HasToy: false, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{HasToy: true, Left: nil, Right: nil},
	}

	fmt.Println(unrollGarland(tree1)) // Ожидаю [false, false, true, true, false]
}
