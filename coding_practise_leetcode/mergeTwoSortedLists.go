package main

import "fmt"

//import "container/list"

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	mergedList := &ListNode{}
	var headOfMergedList *ListNode = mergedList

	for (list1 != nil && list2 != nil) {
		if list1.Val == list2.Val {
			node1 := &ListNode{Val: list1.Val, Next: nil}
			node2 := &ListNode{Val: list2.Val, Next: node1}
			mergedList.Next = node2
			mergedList = mergedList.Next.Next
			list1 = list1.Next
			list2 = list2.Next
		} else if list1.Val < list2.Val {
			node2 := &ListNode{Val: list1.Val}
			mergedList.Next = node2
			mergedList = mergedList.Next
			list1 = list1.Next
		} else if list1.Val > list2.Val {
			node2 := &ListNode{Val: list2.Val}
			mergedList.Next = node2
			mergedList = mergedList.Next
			list2 = list2.Next
		}
	}
	for (list1 != nil){
		if list1 != nil {
			node2 := &ListNode{Val: list1.Val}
			mergedList.Next = node2
			mergedList = mergedList.Next
		}
	}
	for (list2 != nil) {
		if list2 != nil {
			node2 := &ListNode{Val: list2.Val}
			mergedList.Next = node2
			mergedList = mergedList.Next
			list2 = list2.Next
		}
	}
	
	return headOfMergedList.Next
}

func main() {
	node3 := &ListNode{Val: 3, Next: nil}
	node2 := &ListNode{Val: 2, Next: node3}
	node1 := &ListNode{Val: 1, Next: node2}
	
	
	node6 := &ListNode{Val: 4, Next: nil}
	node5 := &ListNode{Val: 3, Next: node6}
	node4 := &ListNode{Val: 1, Next: node5}
	
	mergeTwoLists(node1, node4)
}
