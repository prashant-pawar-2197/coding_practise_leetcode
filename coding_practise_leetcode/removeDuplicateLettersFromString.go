package main

import (
	"fmt"
)

func removeDuplicatesFromStr(s string) string {
	hm := make(map[byte]int)
	lenOfStr := len(s)
	result := ""
	for i := 0; i < lenOfStr; i++ {
		if val, ok := hm[s[i]]; ok {
			hm[s[i]] = val+1
		} else {
			hm[s[i]] = 1
		}
	}
	for i := 0; i < lenOfStr; i++ {
		val, ok := hm[s[i]]
		if val > 1 {
			hm[s[i]] = val
		}
			result += string(s[i])
		}
		if val, ok := hm[s[i]]; ok && val > 1 {
			//skip
			result += string(s[i])
		} else {
			result += string(s[i])
		}
	}
	return result
}

func main(){
	fmt.Println(removeDuplicatesFromStr("Java"))
}