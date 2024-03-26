package main

import (
	"fmt"
)

func strStr(haystack string, needle string) int {
	lenOfHaystack := len(haystack)
	lenOfNeedle := len(needle)
	if lenOfNeedle > lenOfHaystack {
		return -1
	}
    for j:= 0; j<lenOfHaystack; j++ {
		if byte(haystack[j]) == needle[0] {
			dummyJCounter := j
			foundCounter := 1 
			for i := 1; i<lenOfNeedle; i++ {
				dummyJCounter++
				if dummyJCounter >= lenOfHaystack {
					return -1
				}
				if byte(haystack[dummyJCounter]) == needle[i] {
					foundCounter++
					continue
				} else {
					break
				}
			}
			if foundCounter == lenOfNeedle {
				return j
			}
		}
	}
	return -1
}

func main(){
	fmt.Println(strStr("itsSadverySad", "Sad"))
}