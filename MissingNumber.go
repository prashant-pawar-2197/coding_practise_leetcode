package main

import (
	"fmt"
	//"strconv"
)

func missingNumber(nums []int) int {
	maxNum := len(nums)
	reqSum := 0
	currentSum := 0
    for i := 0; i < len(nums); i++ {
		currentSum += nums[i]
	}
	for i := 0; i <= maxNum; i++ {
		reqSum += i	
	}
	return reqSum-currentSum
}

/*
Given an integer array nums, move all 0's to the end of it while maintaining the relative order of the non-zero elements.
Note that you must do this in-place without making a copy of the array.
*/
func moveZeroes(nums []int)  {
    zeroPointer := 0
	nonZeroPointer := 1
	for i := 1; i < len(nums); i++ {
		if nums[zeroPointer] == 0 {
			if nums[nonZeroPointer] != 0 {
				tmpNum := nums[zeroPointer]
				nums[zeroPointer] = nums[nonZeroPointer]
				nums[nonZeroPointer] = tmpNum
				zeroPointer++
				nonZeroPointer++
			} else if nums[nonZeroPointer] == 0 {
				nonZeroPointer++
			}
		} else {
			zeroPointer++
			nonZeroPointer++
		}
	}
}
/*
Given an input string s, reverse the order of the words.

A word is defined as a sequence of non-space characters. The words in s will be separated by at least one space.

Return a string of the words in reverse order concatenated by a single space.

Note that s may contain leading or trailing spaces or multiple spaces between two words. The returned string should only have a single space separating the words. Do not include any extra spaces.

 

Example 1:

Input: s = "the sky is blue"
Output: "blue is sky the"
*/

func reverseWords(s string) string {
	 delimiter := " "
	nameArr := make([]string, 0)
	var word string
	for i := 0; i < len(s); i++ {
		if string(s[i]) != " " {
			word += string(s[i])
		} else {
			if word != "" {
				nameArr = append(nameArr, word)
				word = ""
			}
		}
	}
	if word != "" {
		nameArr = append(nameArr, word)
	}
	finalWord := ""
	for i := len(nameArr)-1; i >=0 ; i-- {
		for _, char := range nameArr[i] {
			finalWord += string(char)
		}
		finalWord += delimiter
	}
	return finalWord
}

func main(){
	var numArr = []int{1,0,1}
	moveZeroes(numArr)
	fmt.Println(numArr)

	var nameArr string = "sky is blue"
	fmt.Println(reverseWords(nameArr))
}