package main


import (
	"fmt"
)

func main(){
	hmm := "T"
	switch hmm {
	case "R": fmt.Println("R")
	default:  fmt.Println("Default")
	}
	fmt.Println("Hi")
}