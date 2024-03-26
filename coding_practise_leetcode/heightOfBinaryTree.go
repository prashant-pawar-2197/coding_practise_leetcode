package main

import (
	"fmt"
)

type TreeNode struct {
    Val int
    Left *TreeNode
    Right *TreeNode
}

func searchElement(root *TreeNode, element int) bool {
	if root == nil {
		return false
	}
	if element == root.Val {
		return true
	}
	if element < root.Val {
		return searchElement(root.Left, element)
	}
	if element > root.Val {
		return searchElement(root.Right, element)
	}
	return false
}

func maxDepth(root *TreeNode) int {	   
	return 1
}

func main(){
	n1 := &TreeNode{Val: 3, Left: nil, Right: nil}
	n3 := &TreeNode{Val: 6, Left: nil, Right: nil}
	n2 := &TreeNode{Val: 5, Left: n1, Right: n3}
	n4 := &TreeNode{Val: 10, Left: nil, Right: nil}
	n5 := &TreeNode{Val: 8, Left: n2, Right: n4}
	n9 := &TreeNode{Val: 15, Left: nil, Right: nil}
	n6 := &TreeNode{Val: 19, Left: nil, Right: nil}
	n8 := &TreeNode{Val: 18, Left: n9, Right: n6}
	n7 := &TreeNode{Val: 12, Left: n5, Right: n8}
	

	fmt.Println(searchElement(n7, 80))
}