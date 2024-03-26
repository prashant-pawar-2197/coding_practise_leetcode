package main

import (
	"fmt"
)
// [1,2,3,0,0,0]
// [2,3,5]

// [3,4,5,6,0,0,0]
// [1,1,2]

func merge(nums1 []int, m int, nums2 []int, n int) {
	p1 := m-1
	if p1 < 0 {
		p1 = 0
	}
	p2 := n-1
	if p2 < 0 {
		p2 = 0
	}
	p := m+n-1
	for(p2 >= 0){
		if nums1[p1] < nums2[p2] {
			nums1[p] = nums2[p2]
			p2--
		} else if nums1[p1] > nums2[p2] {
			nums1[p] = nums1[p1]
			p1--
		} else if nums1[p1] == nums2[p2] {
			nums1[p] = nums2[p2]
			p2--
		}
		p--
	}
	for p2 >= 0 {
		nums1[p] = nums2[p2]
		p--
		p2--
	}
}

func main(){
	var num1 []int = []int{0}
	var num2 []int = []int{1}
	merge(num1, 0, num2, 1)
	fmt.Println(num1)
}