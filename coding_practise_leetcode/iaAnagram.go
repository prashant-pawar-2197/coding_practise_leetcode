package main

import (
	"fmt"
)

func isAnagram(s string, t string) bool {
    lenOfS := len(s)
    lenOfT := len(t)
    if lenOfS != lenOfT {
        return false
    }
    hm := make(map[byte]int)
    //hm1 := make(map[byte]int)
    result := true
    for i:= 0; i < lenOfS; i++ {
        if val, ok := hm[s[i]]; ok {
            hm[s[i]] = val + 1
        } else {
            hm[s[i]] = 1
        }
    }
    for i:= 0; i < lenOfT; i++{
        var (
            val int
            ok bool
        )
        if val, ok = hm[t[i]]; !ok {
            result = false
            break 
        }
        if ok && val == 0 {
            result = false
            break 
        }
        if ok && val >= 1 {
            hm[t[i]] = val - 1
        }
    }
    return result
}

func main(){
	fmt.Println(isAnagram("aacc", "ccac"))
}