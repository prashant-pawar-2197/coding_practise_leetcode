package main

import (
	"fmt"
	"strconv"
)

func main() {
	var (
		dstValue  int
		err       error
		stringVal string
	)
	stringVal = "AAA222"
	if dstValue, err = strconv.Atoi(stringVal); err != nil {
		fmt.Println("Integer conversion failed")
	}
	fmt.Println(dstValue)

}