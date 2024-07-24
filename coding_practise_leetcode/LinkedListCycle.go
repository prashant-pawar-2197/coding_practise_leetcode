/**
 * Definition for singly-linked list.
 */
package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// flyod warshall algorithm (slow and fast pointer)
func hasCycleV2(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	slowPtr := head
	fastPtr := head

	for fastPtr != nil && fastPtr.Next != nil {
		slowPtr = slowPtr.Next
		fastPtr = fastPtr.Next.Next

		if slowPtr == fastPtr {
			return true
		}
	}
	return false
}

// prashant's thought solution
func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	hm := make(map[*ListNode]bool)
	current := head
	result := false
	for current != nil {
		if _, ok := hm[current]; ok {
			result = true
			break
		} else {
			hm[current] = true
		}
		current = current.Next
	}
	return result
}
func main() {
	l2 := &ListNode{Val: 1, Next: nil}
	l1 := &ListNode{Val: 1, Next: l2}
	fmt.Println(hasCycle(l1))
}
