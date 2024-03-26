package main

import "fmt"

//"fmt"

//Definition for a binary tree node.
type TreeNode struct {
    Val int
    Left *TreeNode
    Right *TreeNode
}

func inorderWrapper(root *TreeNode, arr map[int][]int, counter, noOfEle *int)  {
	if root == nil {
		return
	}
	inorderWrapper(root.Left, arr, counter, noOfEle)
	counterArr := make([]int,0)
	var ok bool
	if counterArr, ok = arr[root.Val]; ok {
		counterArr = append(counterArr, *counter)	
	} else {
		counterArr = append(counterArr, *counter)
	}
	arr[root.Val] = counterArr
	*counter++
	*noOfEle++
	inorderWrapper(root.Right, arr, counter, noOfEle)
}

func inorderTraversal(root *TreeNode) []int {
	var counter int = 1
	var noOfEle int = 0
	map1 := make(map[int][]int)
	inorderWrapper(root, map1, &counter, &noOfEle)
	arr := make([]int,noOfEle)
	for k, v := range map1 {
		for _, value := range v {
			arr[value-1] = k			
		}

	}
	return arr
}

func main(){
	
	n1 := &TreeNode{Val: -71, Left: nil, Right: nil}
	n10 := &TreeNode{Val: 8, Left: nil, Right: nil}
	n3 := &TreeNode{Val: -22, Left: nil, Right: n10}
	n2 := &TreeNode{Val: -54, Left: n1, Right: n3}
	n5 := &TreeNode{Val: -100, Left: nil, Right: nil}
	n4 := &TreeNode{Val: -34, Left: nil, Right: n5}
	n7 := &TreeNode{Val: -100, Left: nil, Right: nil}
	n9 := &TreeNode{Val: 48, Left: n2, Right: nil}
	n8 := &TreeNode{Val: -48, Left: n7, Right: n9}
	n6 := &TreeNode{Val: 37, Left: n4, Right: n8}

	fmt.Println(inorderTraversal(n6))


}