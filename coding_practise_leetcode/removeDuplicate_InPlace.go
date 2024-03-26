package main


import (
	"fmt"
)
//[0,0,1,1,1,2,2,3,3,4]
func removeDuplicates(nums []int) int {
	var (
		maxNum, nextPointer, currentPointer, uniqueElems int
	)
	maxNum = nums[0]
	nextPointer = 1
	currentPointer = 0
	lenOfArr := len(nums)

	for nextPointer < lenOfArr {
		if nums[nextPointer] <= maxNum {
			nextPointer++
		} else {
			nums[currentPointer+1] = nums[nextPointer]
			maxNum = nums[currentPointer+1]
			uniqueElems++
			currentPointer++
		}
	}
	fmt.Println(nums)
	fmt.Println(maxNum)
	return uniqueElems+1
}

func main(){
	arr := []int{0,0,0,0,3}
	fmt.Println(removeDuplicates(arr))
}