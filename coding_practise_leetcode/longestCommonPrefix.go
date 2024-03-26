package main

import (
	"fmt"
	"strings"
)

func longestCommonPrefix(strs []string) string {
	smallestStr := 222
	smallestWord := ""
    for _, v := range strs {
		lenOfStr := len(v)
		if lenOfStr < smallestStr {
			smallestStr = lenOfStr
			smallestWord = v
		}
	}
	longestPrefix := smallestWord
	for _, v := range strs {
		found := false
		for v != "" {
			if strings.Contains(v, longestPrefix) {
				found = true
				break
			} else {
				longestPrefix = longestPrefix[:len(longestPrefix)-1]
			}
		}
		if !found {
			return ""
		}
	}
	return longestPrefix
}

func main(){
	var strArr = []string{"flower", "flow", "flight"}
	fmt.Println(longestCommonPrefix(strArr))
}