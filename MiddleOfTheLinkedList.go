/*
	Given the head of a singly linked list, return the middle node of the linked list.
	If there are two middle nodes, return the second middle node.
*/

package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func middleNode(head *ListNode) *ListNode {
	var fastPointer *ListNode = head
	var slowPointer *ListNode = head

	for fastPointer != nil && fastPointer.Next != nil {
		fastPointer = fastPointer.Next.Next
		slowPointer = slowPointer.Next
	}
	return slowPointer
}

func main() {

}
