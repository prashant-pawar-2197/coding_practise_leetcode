package main

import (
	"container/list"
	"fmt"
)

type n struct {
	Val int
	prev *List
	next *List
}

type DoublyLinkedList struct {
	head *List
	tail *List
}

func plusOne(digits []int) []int {
	lenOfDigits := len(digits)
	carryForward := 0
	isCarryForward := false
	newArr := make
	doublyList := &DoublyLinkedList{head: nil, tail: nil}
	for i:=lenOfDigits-1; i >= 0 ; i++ {
		if digits[i] == 9 {
			isCarryForward = true
			carryForward = 1
			list
		}
	}
}

func main(){

}