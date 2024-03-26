package main

import (
	"fmt"
	"strconv"
	"strings"
)

func wordPattern(pattern string, s string) bool {
    hm := make(map[byte]int64)
    hm2 := make(map[string]int64) 
    var counter int64 = 1
    for i:= 0; i< len(pattern); i++ {
        if _, ok := hm[pattern[i]]; !ok {
            hm[pattern[i]] = counter
            counter++
        }
    }
    numPat1 := ""
    numPat2 := ""

    for i:= 1; i< len(pattern); i++ {
        if val, ok := hm[pattern[i]]; ok {
           numPat1 += strconv.FormatInt(val, 10)
        }
    }
    strArr := strings.Split(s, " ")
    counter = 1
    for i:= 0; i< len(strArr); i++ {
        if _, ok := hm2[strArr[i]]; !ok {
            hm2[strArr[i]] = counter
            counter++
        }
    }

    for i:= 1; i< len(strArr); i++ {
        if val, ok := hm2[strArr[i]]; ok {
           numPat2 += strconv.FormatInt(val, 10)
        }
    }

    if numPat1 == numPat2 {
        return true
    }
    return false
}

func main(){
	fmt.Println(wordPattern("abba", "dog cat cat dog"))
}