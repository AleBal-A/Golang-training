package main

import (
	"reflect"
	"testing"
)

func TestUnrollGarland(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []bool
	}{
		{
			name: "Test case 1",
			root: &TreeNode{
				HasToy: false,
				Left: &TreeNode{
					HasToy: false,
					Left:   &TreeNode{HasToy: false, Left: nil, Right: nil},
					Right:  &TreeNode{HasToy: true, Left: nil, Right: nil},
				},
				Right: &TreeNode{HasToy: true, Left: nil, Right: nil},
			},
			expected: []bool{false, false, true, true, false},
		},
		{
			name:     "Empty tree",
			root:     nil,
			expected: []bool{},
		},
		{
			name: "Single node",
			root: &TreeNode{
				HasToy: true,
			},
			expected: []bool{true},
		},
		{
			name: "Two levels",
			root: &TreeNode{
				HasToy: true,
				Left:   &TreeNode{HasToy: false},
				Right:  &TreeNode{HasToy: true},
			},
			expected: []bool{true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unrollGarland(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unrollGarland(%v) = %v; expected %v", tt.root, result, tt.expected)
			}
		})
	}
}
