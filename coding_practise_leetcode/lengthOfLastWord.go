package main

import (
	"fmt"
)

/*
Given a string s consisting of words and spaces, return the length of the last word in the string.
A word is a maximal substring consisting of non-space characters only.
*/

func lengthOfLastWord(s string) int {
    counter := 0
    lengthSlice := make([]int,0)
    for i:=0 ; i < len(s); i++ {
        if string(s[i]) != " " {
            counter++
        } else {
            lengthSlice = append(lengthSlice, counter)
            counter = 0
        }
    }
	if counter != 0 {
		lengthSlice = append(lengthSlice, counter)
        counter = 0
	}
	result := 0
    lenOfSlice := len(lengthSlice)
    if lenOfSlice > 0 {
		for lenOfSlice >= 0 {
			if lengthSlice[lenOfSlice-1] != 0 {
				result = lengthSlice[lenOfSlice-1]
				break
			}
			lenOfSlice--
		}
    }
    return result
}

func main(){
	fmt.Println(lengthOfLastWord("luffy is still joyboy"))
}