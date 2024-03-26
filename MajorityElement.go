/*
Given an array nums of size n, return the majority element.
The majority element is the element that appears more than ⌊n / 2⌋ times. 
You may assume that the majority element always exists in the array.

Example 1:
Input: nums = [3,2,3]
Output: 3

Example 2:
Input: nums = [2,2,1,1,1,2,2]
Output: 2
*/

package main

import (
	"fmt"
)

func majorityElement(nums []int) int {
	numOccurences := make(map[int]int)
    for _, v := range nums {
		if val, ok := numOccurences[v]; ok {
			numOccurences[v] = val+1
		} else {
			numOccurences[v] = 1
		}
	}
	maxRep := len(nums)/2
	for k, v := range numOccurences {
		if v > maxRep {
			return k
		}
	}
	return 0
}

func main(){

}