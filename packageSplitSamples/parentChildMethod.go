package main

import (
	"fmt"
)
//var multiplier int

func mathematicsOld(a , b int){
	multiplier := 2
	muliply(a,b)
	multiply := func (a, b int) int {
		return a*multiplier*b
	}
}	



func main()  {
	mathematicsOld(10,2)
}