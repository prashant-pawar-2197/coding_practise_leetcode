package main

import (
	"fmt"
)

func mySqrt(x int) int {
	if x == 0 || x == 1 {
		return x
	} else {
		left := 2
		right := x/2
		result := 0
		for left <= right {
			mid := left + (right-left)/2
			sqOfMid := mid*mid
			if sqOfMid == x {
				return mid
			} else if sqOfMid < x {
				left = mid
				result = mid
			} else {
				right = mid
			}
		}
		return result
	}
}

func main(){
	fmt.Println(mySqrt(8))
}