package main

import (
	"fmt"
	//"math"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func hasPathSum(root *TreeNode, targetSum int) bool {
	var sum int
	return hasPathSumWrapper(root, sum, targetSum)
}

func hasPathSumWrapper(root *TreeNode, sum int, targetSum int) bool {
	if root == nil {

		return false

	}
	sum += root.Val

	if root.Left == nil && root.Right == nil {
		if sum == targetSum {
			return true
		} else {
			return false
		}
	}
	if root.Left != nil {
		if !hasPathSumWrapper(root.Left, sum, targetSum) {
			if root.Right != nil {
				return hasPathSumWrapper(root.Right, sum, targetSum)
			}
		} else {
			return true
		}
	} else if root.Right != nil {
		return hasPathSumWrapper(root.Right, sum, targetSum)
	}
	return false
}

func main() {
	/*
	n9 := &TreeNode{Val: 2, Left: nil, Right: nil}
	n8 := &TreeNode{Val: 7, Left: nil, Right: nil}
	n7 := &TreeNode{Val: 1, Left: nil, Right: nil}
	n6 := &TreeNode{Val: 4, Left: nil, Right: n7}
	n5 := &TreeNode{Val: 13, Left: nil, Right: nil}
	n4 := &TreeNode{Val: 11, Left: n8, Right: n9}
	n3 := &TreeNode{Val: 8, Left: n5, Right: n6}
	n2 := &TreeNode{Val: 4, Left: n4, Right: nil}
	n1 := &TreeNode{Val: 5, Left: n2, Right: n3}

*/
	//new node for case
	n12 := &TreeNode{Val: 3, Left: nil, Right: nil}
	n11 := &TreeNode{Val: -2, Left: nil, Right: nil}
	n10 := &TreeNode{Val: 1, Left: n11, Right: n12}
	fmt.Println(hasPathSum(n10, 3))
}