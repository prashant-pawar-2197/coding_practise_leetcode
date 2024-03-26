package main

import (
	"fmt"
)
// 3, 2, 2, 3
func removeElement(nums []int, val int) int {
	var (
	deletePointer int = 0
	shiftPointer  int = 1
)
lenOfArr := len(nums)
if ((lenOfArr == 1 && nums[deletePointer] == val) || lenOfArr == 0) {
	return 0
}
if (lenOfArr == 1 && nums[deletePointer] != val){
	return 1
}
for shiftPointer < lenOfArr {
	if nums[deletePointer] == val {
		tmp := nums[deletePointer]
		if nums[shiftPointer] == tmp {
			shiftPointer++
		} else {
			nums[deletePointer] = nums[shiftPointer]
			nums[shiftPointer] = tmp
			shiftPointer++
			deletePointer++
		}
	} else {
		shiftPointer++
		deletePointer++
	}
}
var result int
for i := 0 ; i < lenOfArr; i++{
	if nums[i] == val {
		result = i
		break
	}
}
if result == 0 && lenOfArr != 1{
	return lenOfArr
}
return result
}

func main(){
	var arr = []int{3,3}
	fmt.Println(removeElement(arr,3))
}