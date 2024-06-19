package main

import (
	"testing"
)

func TestAreToysBalanced(t *testing.T) {
	// Test case 1: Balanced tree
	tree1 := &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: false,
			Left:   &TreeNode{HasToy: false, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{HasToy: true, Left: nil, Right: nil},
	}

	if !areToysBalanced(tree1) {
		t.Errorf("Expected true but got false")
	}

	// Test case 2: Unbalanced tree
	tree2 := &TreeNode{
		HasToy: true,
		Left: &TreeNode{
			HasToy: true,
			Left:   nil,
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{HasToy: false, Left: nil, Right: nil},
	}

	if areToysBalanced(tree2) {
		t.Errorf("Expected false but got true")
	}

	// This tree is unbalanced (3 toys on the left, 2 toys on the right)
	tree3 := &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: true,
			Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{
			HasToy: false,
			Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: false, Left: nil, Right: &TreeNode{HasToy: true}},
		},
	}

	if areToysBalanced(tree3) {
		t.Errorf("Expected false but got true")
	}

	// This tree is balanced (2 toys on the left, 2 toys on the right)
	tree4 := &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: false,
			Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{
			HasToy: false,
			Left:   &TreeNode{HasToy: true, Left: nil, Right: nil},
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
	}

	if !areToysBalanced(tree4) {
		t.Errorf("Expected true but got false")
	}
}

func TestCountToys(t *testing.T) {
	tree := &TreeNode{
		HasToy: true,
		Left: &TreeNode{
			HasToy: true,
			Left:   nil,
			Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
		},
		Right: &TreeNode{HasToy: false, Left: nil, Right: nil},
	}

	expectedCount := 3
	if countToys(tree) != expectedCount {
		t.Errorf("Expected %d but got %d", expectedCount, countToys(tree))
	}
}
