package main


import (
	"fmt"
)

func rearrangeArray(nums []int) []int {
    numPositive := make([]int,0)
    numNegative := make([]int,0)
    lenOfArr := len(nums)
    for i:=0; i < lenOfArr; i++ {
        if nums[i] < 0 {
            numNegative = append(numNegative, nums[i])
        } else {
            numPositive = append(numPositive, nums[i])
        }
    }
    lenOfNeg := len(numNegative)
    lenOfPos := len(numPositive)
    finalArrLength := len(numPositive) + len(numNegative)
    finalArr := make([]int,0)
    low := 0
    high := 0
    for i:=0; i<finalArrLength; i++{
        if low < lenOfPos && i%2 == 0 {
            finalArr = append(finalArr, numPositive[low])
            low++
			continue
        }
        if high < lenOfNeg && i%2 == 1 {
            finalArr = append(finalArr, numNegative[high])
            high++
			continue
        }
    }
    return finalArr
}

func main(){
	arr := []int{3,1,-2,-5,2,-4}
	fmt.Println(rearrangeArray(arr))
}