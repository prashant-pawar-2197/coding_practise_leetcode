package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var x int
	size := unsafe.Sizeof(x)
	fmt.Printf("Size of int on this system: %d bytes\n", size)
}
