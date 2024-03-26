package main

import (
	"fmt"
)

func main(){
	listOfEmptyFields := ""
	delimiter := ","
	arr := [4]string{"AAA", "BBB", "CCC", "DDD"}
	cdr := make(map[string]interface{})
	for i := 0; i < 4; i++ {
		if cdr[arr[i]] == nil {
			listOfEmptyFields = listOfEmptyFields + delimiter + arr[i]
		}
	}
	fmt.Println(listOfEmptyFields)
}