package main 

import (
	"fmt"
)

func isPalindrome(s string) bool {
    lenOfStr := len(s)
    low := 0
    high := lenOfStr-1
    result := true
    for low < high {
        lowChar := int(s[low])
        highChar := int(s[high])
        if (lowChar >= 33 && lowChar <= 47) || (lowChar >= 58 && lowChar <= 64) || (lowChar >= 91 && lowChar <= 96) || (lowChar >= 123 && lowChar <= 127) {
            low++
            continue
        }
        if (highChar >= 33 && highChar <= 47) || (highChar >= 58 && highChar <= 64) || (highChar >= 91 && highChar <= 96) || (highChar >= 123 && highChar <= 127) {
            high--
            continue
        }
        if lowChar >= 65 && lowChar <= 90 {
            lowChar += 32
        }
        if highChar >= 65 && highChar <= 90 {
            highChar += 32
        }
        if lowChar != highChar {
            result = false
            break
        }
        low++
        high--
    }
    return result
}


func main(){
	fmt.Println(isPalindrome("A man, a plan, a canal: Panama"))
}