package main

import (
	"fmt"
)

type ListNode struct {
     Val int
     Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
    if head.Next == nil {
        return head
    }
    current := head
    next := head.Next
    for next != nil {
        if current == next {
			if next.Next != nil {
				next = next.Next
			}
        }
        if current.Val == next.Val {
            next = next.Next
            current.Next = next
        } else {
            current = next
            next = next.Next
        }
    }
    return head
}

func main(){
	l5 := &ListNode{Val: 3, Next: nil}
	l4 := &ListNode{Val: 3, Next: l5}
	l3 := &ListNode{Val: 2, Next: l4}
	l2 := &ListNode{Val: 2, Next: l3}
	l1 := &ListNode{Val: 1, Next: l2}

	l11 := deleteDuplicates(l1)
	for l11 != nil {
		fmt.Println(l11.Val)
		l11 = l11.Next
	} 
}